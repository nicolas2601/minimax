package recurring

import (
	"time"

	"github.com/google/uuid"
)

// Frequency determines how often a recurring rule fires.
type Frequency string

const (
	FrequencyDaily   Frequency = "daily"
	FrequencyWeekly  Frequency = "weekly"
	FrequencyBiweekly Frequency = "biweekly"
	FrequencyMonthly Frequency = "monthly"
	FrequencyYearly  Frequency = "yearly"
)

func IsValidFrequency(f string) bool {
	switch Frequency(f) {
	case FrequencyDaily, FrequencyWeekly, FrequencyBiweekly, FrequencyMonthly, FrequencyYearly:
		return true
	}
	return false
}

// TxType mirrors transactions.Type but lives in this package to avoid an
// import dependency on transactions. Validation happens against this enum
// when the rule is created.
type TxType string

const (
	TypeExpense TxType = "expense"
	TypeIncome  TxType = "income"
)

func IsValidTxType(t string) bool {
	switch TxType(t) {
	case TypeExpense, TypeIncome:
		return true
	}
	return false
}

type RunStatus string

const (
	RunPending  RunStatus = "pending"
	RunExecuted RunStatus = "executed"
	RunSkipped  RunStatus = "skipped"
	RunFailed   RunStatus = "failed"
)

type Rule struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	AccountID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"account_id"`
	CategoryID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"category_id"`
	Type         TxType     `gorm:"not null;size:20" json:"type"`
	Amount       int64      `gorm:"not null" json:"amount"`
	Currency     string     `gorm:"not null;size:3;default:COP" json:"currency"`
	Description  *string    `gorm:"size:255" json:"description,omitempty"`
	Notes        *string    `gorm:"type:text" json:"notes,omitempty"`
	Frequency    Frequency  `gorm:"not null;size:20" json:"frequency"`
	IntervalCount int       `gorm:"not null;default:1" json:"interval_count"`
	StartDate    time.Time  `gorm:"not null" json:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	LastRunDate  *time.Time `json:"last_run_date,omitempty"`
	NextRunDate  time.Time  `gorm:"not null;index" json:"next_run_date"`
	IsActive     bool       `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Rule) TableName() string { return "recurring_rules" }

// Run records each (scheduled or actual) execution attempt. UNIQUE(rule_id,
// scheduled_date) makes generation idempotent — a re-run of generate-today
// won't duplicate transactions.
type Run struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	RecurringRuleID uuid.UUID `gorm:"type:uuid;not null;index" json:"recurring_rule_id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	ScheduledDate  time.Time  `gorm:"not null" json:"scheduled_date"`
	ExecutedAt     *time.Time `json:"executed_at,omitempty"`
	Status         RunStatus  `gorm:"not null;size:20" json:"status"`
	TransactionID  *uuid.UUID `gorm:"type:uuid;index" json:"transaction_id,omitempty"`
	ErrorMessage   *string    `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (Run) TableName() string { return "recurring_runs" }

// NextOccurrence returns the next scheduled date strictly after `from` given
// the rule's frequency & interval. EndDate-inclusive (returns nil if the
// next date would be after EndDate). `advance` always returns a valid
// time.Time for any cadence the model understands, so the only "no
// occurrence" signal is the EndDate guard below.
func (r *Rule) NextOccurrence(from time.Time) *time.Time {
	candidate := advance(from, r.StartDate, r.Frequency, r.IntervalCount, 1)
	if r.EndDate != nil && candidate.After(*r.EndDate) {
		return nil
	}
	return &candidate
}

func (r *Rule) OccurrencesBetween(from, to time.Time) []time.Time {
	if to.Before(from) {
		return nil
	}
	out := []time.Time{}
	current := from
	for {
		next := r.NextOccurrence(current)
		if next == nil || next.After(to) {
			return out
		}
		out = append(out, *next)
		current = *next
	}
}

// advance advances start by k cycles of (frequency * interval_count). It
// returns the first date strictly after `from` that lines up with the
// cadence. Pure function — no DB / clock side-effects.
func advance(from, start time.Time, f Frequency, interval, k int) time.Time {
	// Walk k steps at a time. For small interval_count this is fine; for
	// pathological cases (interval=100, monthly) we'd want a closed-form
	// jump, but for personal-finance use it's negligible.
	cursor := start
	for step := 0; step < k; step++ {
		cursor = stepForward(cursor, f, interval)
	}
	// If we haven't passed `from` yet, walk forward until we do.
	for !cursor.After(from) {
		cursor = stepForward(cursor, f, interval)
	}
	return cursor
}

func stepForward(t time.Time, f Frequency, interval int) time.Time {
	switch f {
	case FrequencyDaily:
		return t.AddDate(0, 0, interval)
	case FrequencyWeekly:
		return t.AddDate(0, 0, 7*interval)
	case FrequencyBiweekly:
		return t.AddDate(0, 0, 14*interval)
	case FrequencyMonthly:
		return t.AddDate(0, interval, 0)
	case FrequencyYearly:
		return t.AddDate(interval, 0, 0)
	}
	return t
}
