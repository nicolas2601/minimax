package recurring

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidAmount    = errors.New("amount must be greater than zero")
	ErrInvalidFrequency = errors.New("frequency must be one of daily/weekly/biweekly/monthly/yearly")
	ErrInvalidType      = errors.New("type must be expense or income")
	ErrInvalidInterval  = errors.New("interval_count must be >= 1")
	ErrInvalidStartDate = errors.New("start_date is required")
	ErrEndBeforeStart   = errors.New("end_date must be on or after start_date")
)

// AccountLookup + CategoryLookup are the minimal contracts recurring needs
// to validate that account_id and category_id exist and belong to the user.
type AccountLookup interface {
	Exists(id, userID uuid.UUID) (bool, error)
}

type CategoryLookup interface {
	Exists(id, userID uuid.UUID) (bool, error)
}

// TxCreator is the contract for creating one transaction. The concrete
// implementation lives in transactions.Service via an adapter in main.go
// to keep this package independent.
type TxCreator interface {
	CreateFromRecurring(userID, accountID, categoryID uuid.UUID, txType string,
		amount int64, currency string, date time.Time, description, notes *string,
		recurringRunID uuid.UUID) (uuid.UUID, error)
}

// UserResolver resolves a userID string (as stored on rules) back to a
// uuid.UUID for run insertion. Returns ErrRuleNotFound if the user is gone
// (shouldn't happen with FK cascade but defensive).
type UserResolver interface {
	UserIDFromString(userIDStr string) (uuid.UUID, error)
}

type Service struct {
	repo        Repository
	accounts    AccountLookup
	categories  CategoryLookup
	txCreator   TxCreator
	users       UserResolver
	now         func() time.Time // injectable for tests
}

func NewService(repo Repository, accounts AccountLookup, categories CategoryLookup, txCreator TxCreator, users UserResolver) *Service {
	return &Service{
		repo:       repo,
		accounts:   accounts,
		categories: categories,
		txCreator:  txCreator,
		users:      users,
		now:        time.Now,
	}
}

type CreateRequest struct {
	AccountID     string  `json:"account_id"`
	CategoryID    string  `json:"category_id"`
	Type          string  `json:"type"`
	Amount        int64   `json:"amount"`
	Currency      string  `json:"currency"`
	Description   *string `json:"description,omitempty"`
	Notes         *string `json:"notes,omitempty"`
	Frequency     string  `json:"frequency"`
	IntervalCount int     `json:"interval_count"`
	StartDate     string  `json:"start_date"`
	EndDate       *string `json:"end_date,omitempty"`
}

type UpdateRequest struct {
	Amount        *int64  `json:"amount,omitempty"`
	Description   *string `json:"description,omitempty"`
	Notes         *string `json:"notes,omitempty"`
	Frequency     *string `json:"frequency,omitempty"`
	IntervalCount *int    `json:"interval_count,omitempty"`
	EndDate       *string `json:"end_date,omitempty"`
	ClearEndDate  bool    `json:"clear_end_date,omitempty"`
	IsActive      *bool   `json:"is_active,omitempty"`
}

func (s *Service) Create(userID uuid.UUID, req CreateRequest) (*Rule, error) {
	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if !IsValidFrequency(req.Frequency) {
		return nil, ErrInvalidFrequency
	}
	if !IsValidTxType(req.Type) {
		return nil, ErrInvalidType
	}
	interval := req.IntervalCount
	if interval < 1 {
		interval = 1
	}
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, ErrInvalidStartDate
	}
	var end *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		e, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, ErrInvalidStartDate
		}
		if e.Before(start) {
			return nil, ErrEndBeforeStart
		}
		end = &e
	}
	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		return nil, fmt.Errorf("parse account_id: %w", err)
	}
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("parse category_id: %w", err)
	}
	if s.accounts != nil {
		ok, err := s.accounts.Exists(accountID, userID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("account not found or not owned by user")
		}
	}
	if s.categories != nil {
		ok, err := s.categories.Exists(categoryID, userID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("category not found or not owned by user")
		}
	}
	currency := req.Currency
	if currency == "" {
		currency = "COP"
	}

	// First occurrence must be strictly after start (or = start if today >= start).
	next := start
	today := s.now().Truncate(24 * time.Hour)
	if !next.After(today) {
		// Find first strictly-after today.
		cursor := start
		for !cursor.After(today) {
			cursor = stepForward(cursor, Frequency(req.Frequency), interval)
		}
		next = cursor
	}

	rr := &Rule{
		UserID:        userID,
		AccountID:     accountID,
		CategoryID:    categoryID,
		Type:          TxType(req.Type),
		Amount:        req.Amount,
		Currency:      currency,
		Description:   req.Description,
		Notes:         req.Notes,
		Frequency:     Frequency(req.Frequency),
		IntervalCount: interval,
		StartDate:     start,
		EndDate:       end,
		NextRunDate:   next,
		IsActive:      true,
	}
	if err := s.repo.CreateRule(rr); err != nil {
		return nil, err
	}
	return rr, nil
}

func (s *Service) Get(id, userID uuid.UUID) (*Rule, error) {
	return s.repo.GetRuleByID(id, userID)
}

func (s *Service) List(userID uuid.UUID) ([]Rule, error) {
	return s.repo.ListRulesByUser(userID)
}

func (s *Service) Update(id, userID uuid.UUID, req UpdateRequest) (*Rule, error) {
	rr, err := s.repo.GetRuleByID(id, userID)
	if err != nil {
		return nil, err
	}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return nil, ErrInvalidAmount
		}
		rr.Amount = *req.Amount
	}
	if req.Description != nil {
		rr.Description = req.Description
	}
	if req.Notes != nil {
		rr.Notes = req.Notes
	}
	if req.Frequency != nil {
		if !IsValidFrequency(*req.Frequency) {
			return nil, ErrInvalidFrequency
		}
		rr.Frequency = Frequency(*req.Frequency)
	}
	if req.IntervalCount != nil {
		if *req.IntervalCount < 1 {
			return nil, ErrInvalidInterval
		}
		rr.IntervalCount = *req.IntervalCount
	}
	if req.ClearEndDate {
		rr.EndDate = nil
	} else if req.EndDate != nil && *req.EndDate != "" {
		e, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, ErrInvalidStartDate
		}
		if e.Before(rr.StartDate) {
			return nil, ErrEndBeforeStart
		}
		rr.EndDate = &e
	}
	if req.IsActive != nil {
		rr.IsActive = *req.IsActive
	}
	if err := s.repo.UpdateRule(rr); err != nil {
		return nil, err
	}
	return rr, nil
}

func (s *Service) Delete(id, userID uuid.UUID) error {
	return s.repo.DeleteRule(id, userID)
}

// GenerateToday processes every active rule whose next_run_date <= today.
// Idempotent: re-running the same day won't create duplicate transactions
// because (rule_id, scheduled_date) is UNIQUE in recurring_runs.
func (s *Service) GenerateToday() (stats GenerateStats, err error) {
	today := s.now().Truncate(24 * time.Hour)
	rules, err := s.repo.ListDueRules(today)
	if err != nil {
		return stats, err
	}
	for _, rule := range rules {
		stats.RulesScanned++
		nextOcc := rule.NextOccurrence(rule.LastRunDateOrZero())
		if nextOcc == nil {
			// No more occurrences (e.g. past EndDate) — deactivate.
			rule.IsActive = false
			_ = s.repo.UpdateRule(&rule)
			stats.RulesDeactivated++
			continue
		}
		if nextOcc.After(today) {
			continue
		}
		// Walk through all occurrences up to today.
		for nextOcc != nil && !nextOcc.After(today) {
			// Skip if already executed.
			if existing, _ := s.repo.GetRunByRuleAndDate(rule.ID, *nextOcc); existing != nil {
				nextOcc = rule.NextOccurrence(*nextOcc)
				continue
			}
			run := Run{
				ID:             uuid.New(),
				RecurringRuleID: rule.ID,
				UserID:         rule.UserID,
				ScheduledDate:  *nextOcc,
				Status:         RunPending,
			}
			if err := s.repo.CreateRun(&run); err != nil {
				stats.Errors = append(stats.Errors, fmt.Errorf("rule %s on %s: create run: %w", rule.ID, nextOcc.Format("2006-01-02"), err))
				nextOcc = rule.NextOccurrence(*nextOcc)
				continue
			}
			if s.txCreator == nil {
				run.Status = RunSkipped
				run.ErrorMessage = stringPtr("transactions creator not configured")
				_ = s.repo.UpdateRun(&run)
				stats.Errors = append(stats.Errors, fmt.Errorf("rule %s: no tx creator", rule.ID))
				nextOcc = rule.NextOccurrence(*nextOcc)
				continue
			}
			txID, err := s.txCreator.CreateFromRecurring(
				rule.UserID, rule.AccountID, rule.CategoryID,
				string(rule.Type), rule.Amount, rule.Currency,
				*nextOcc, rule.Description, rule.Notes, run.ID,
			)
			now := s.now()
			if err != nil {
				run.Status = RunFailed
				msg := err.Error()
				run.ErrorMessage = &msg
				_ = s.repo.UpdateRun(&run)
				stats.Errors = append(stats.Errors, fmt.Errorf("rule %s on %s: %w", rule.ID, nextOcc.Format("2006-01-02"), err))
			} else {
				run.Status = RunExecuted
				run.ExecutedAt = &now
				run.TransactionID = &txID
				_ = s.repo.UpdateRun(&run)
				stats.TransactionsCreated++
			}
			last := *nextOcc
			rule.LastRunDate = &last
			nextOcc = rule.NextOccurrence(last)
		}
		// Persist last_run_date and bump next_run_date.
		rule.NextRunDate = nextOccOrToday(rule, today)
		if err := s.repo.UpdateRule(&rule); err != nil {
			stats.Errors = append(stats.Errors, fmt.Errorf("rule %s: update after generate: %w", rule.ID, err))
		}
	}
	return stats, nil
}

// RunNow executes a single rule immediately, regardless of schedule.
// Creates one transaction for "today" and bumps next_run_date to the next
// future occurrence.
func (s *Service) RunNow(id, userID uuid.UUID) (uuid.UUID, error) {
	rule, err := s.repo.GetRuleByID(id, userID)
	if err != nil {
		return uuid.Nil, err
	}
	if !rule.IsActive {
		return uuid.Nil, fmt.Errorf("rule is not active")
	}
	if s.txCreator == nil {
		return uuid.Nil, fmt.Errorf("transactions creator not configured")
	}
	today := s.now().Truncate(24 * time.Hour)
	run := Run{
		ID:             uuid.New(),
		RecurringRuleID: rule.ID,
		UserID:         rule.UserID,
		ScheduledDate:  today,
		Status:         RunPending,
	}
	if err := s.repo.CreateRun(&run); err != nil {
		return uuid.Nil, err
	}
	txID, err := s.txCreator.CreateFromRecurring(
		rule.UserID, rule.AccountID, rule.CategoryID,
		string(rule.Type), rule.Amount, rule.Currency,
		today, rule.Description, rule.Notes, run.ID,
	)
	now := s.now()
	if err != nil {
		run.Status = RunFailed
		msg := err.Error()
		run.ErrorMessage = &msg
		_ = s.repo.UpdateRun(&run)
		return uuid.Nil, err
	}
	run.Status = RunExecuted
	run.ExecutedAt = &now
	run.TransactionID = &txID
	_ = s.repo.UpdateRun(&run)
	rule.LastRunDate = &today
	rule.NextRunDate = nextOccOrToday(*rule, today)
	_ = s.repo.UpdateRule(rule)
	return txID, nil
}

// ListRuns returns the N most recent runs for a rule, newest first.
func (s *Service) ListRuns(ruleID, userID uuid.UUID, limit int) ([]Run, error) {
	if _, err := s.repo.GetRuleByID(ruleID, userID); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return s.repo.ListRunsByRule(ruleID, limit)
}

type GenerateStats struct {
	RulesScanned        int      `json:"rules_scanned"`
	RulesDeactivated    int      `json:"rules_deactivated"`
	TransactionsCreated int      `json:"transactions_created"`
	Errors              []error  `json:"-"`
	ErrorMessages       []string `json:"errors,omitempty"`
}

// AsDTO exposes the stats shape to JSON callers (Gin can't serialize
// []error directly).
func (g GenerateStats) AsDTO() map[string]any {
	out := map[string]any{
		"rules_scanned":         g.RulesScanned,
		"rules_deactivated":     g.RulesDeactivated,
		"transactions_created":  g.TransactionsCreated,
	}
	if len(g.Errors) > 0 {
		msgs := make([]string, 0, len(g.Errors))
		for _, e := range g.Errors {
			msgs = append(msgs, e.Error())
		}
		out["errors"] = msgs
	}
	return out
}

// --- helpers ---

func (r *Rule) LastRunDateOrZero() time.Time {
	if r.LastRunDate != nil {
		return *r.LastRunDate
	}
	return r.StartDate.AddDate(0, 0, -1) // strictly before start so NextOccurrence picks start
}

func nextOccOrToday(rule Rule, today time.Time) time.Time {
	n := rule.NextOccurrence(rule.LastRunDateOrZero())
	if n == nil {
		return today
	}
	return *n
}

func stringPtr(s string) *string { return &s }