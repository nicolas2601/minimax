package travel

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrGroupNotFound     = errors.New("travel group not found")
	ErrExpenseNotFound   = errors.New("travel expense not found")
	ErrMemberNotFound    = errors.New("member not found in group")
	ErrAlreadyMember     = errors.New("user is already a member of this group")
	ErrNotMember         = errors.New("user is not a member of this group")
	ErrInvalidSplit      = errors.New("invalid split method")
	ErrSplitSumMismatch  = errors.New("split amounts must sum to expense total")
	ErrSplitUsersMissing = errors.New("at least one user must be included in the split")
	ErrPayerNotMember    = errors.New("payer must be a member of the group")
	ErrCannotRemoveOwner = errors.New("cannot remove the group owner")
)

type Repository interface {
	// Groups
	CreateGroup(g *TravelGroup) error
	GetGroup(id uuid.UUID) (*TravelGroup, error)
	ListGroupsByUser(userID uuid.UUID) ([]TravelGroup, error)
	UpdateGroup(g *TravelGroup) error
	DeleteGroup(id uuid.UUID) error

	// Members
	AddMember(m *TravelGroupMember) error
	ListMembers(groupID uuid.UUID) ([]TravelGroupMember, error)
	GetMember(groupID, userID uuid.UUID) (*TravelGroupMember, error)
	RemoveMember(groupID, userID uuid.UUID) error
	CountOwners(groupID uuid.UUID) (int, error)

	// Expenses
	CreateExpenseWithShares(expense *TravelExpense, shares []TravelExpenseShare) error
	GetExpense(id uuid.UUID) (*TravelExpense, error)
	ListExpensesByGroup(groupID uuid.UUID) ([]TravelExpense, error)
	DeleteExpense(id uuid.UUID) error
	ListSharesByExpense(expenseID uuid.UUID) ([]TravelExpenseShare, error)

	// Settlements
	CreateSettlement(s *TravelSettlement) error
	GetSettlement(id uuid.UUID) (*TravelSettlement, error)
	ListSettlementsByGroup(groupID uuid.UUID) ([]TravelSettlement, error)
	UpdateSettlement(s *TravelSettlement) error

	// Aggregations
	SumPaidByUser(groupID, userID uuid.UUID) (int64, error)
	SumShareByUser(groupID, userID uuid.UUID) (int64, error)
	BalancesByGroup(groupID uuid.UUID) ([]MemberBalance, error)
	ListExpensesWithSharesByGroup(groupID uuid.UUID) ([]ExpenseWithShares, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

// --- Groups ---

func (r *repo) CreateGroup(g *TravelGroup) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return r.db.Create(g).Error
}

func (r *repo) GetGroup(id uuid.UUID) (*TravelGroup, error) {
	var g TravelGroup
	err := r.db.First(&g, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrGroupNotFound
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *repo) ListGroupsByUser(userID uuid.UUID) ([]TravelGroup, error) {
	var groups []TravelGroup
	err := r.db.
		Joins("JOIN travel_group_members m ON m.group_id = travel_groups.id").
		Where("m.user_id = ?", userID).
		Order("travel_groups.created_at DESC").
		Find(&groups).Error
	return groups, err
}

func (r *repo) UpdateGroup(g *TravelGroup) error {
	return r.db.Save(g).Error
}

func (r *repo) DeleteGroup(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&TravelGroup{}).Error
}

// --- Members ---

func (r *repo) AddMember(m *TravelGroupMember) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.JoinedAt.IsZero() {
		m.JoinedAt = time.Now()
	}
	// ON CONFLICT DO NOTHING + RETURNING converts the UNIQUE constraint
	// race into a clean no-op. We then re-fetch if RETURNING returned
	// 0 rows (handled in the service after this returns nil).
	res := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(m)
	if res.Error != nil {
		// Some drivers return ErrDuplicatedKey on conflict; tolerate it.
		if isUniqueViolation(res.Error) {
			return ErrAlreadyMember
		}
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrAlreadyMember
	}
	return nil
}

func (r *repo) ListMembers(groupID uuid.UUID) ([]TravelGroupMember, error) {
	var ms []TravelGroupMember
	err := r.db.Where("group_id = ?", groupID).
		Order("joined_at ASC").
		Find(&ms).Error
	return ms, err
}

func (r *repo) GetMember(groupID, userID uuid.UUID) (*TravelGroupMember, error) {
	var m TravelGroupMember
	err := r.db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrMemberNotFound
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repo) RemoveMember(groupID, userID uuid.UUID) error {
	return r.db.Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&TravelGroupMember{}).Error
}

func (r *repo) CountOwners(groupID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Model(&TravelGroupMember{}).
		Where("group_id = ? AND role = ?", groupID, RoleOwner).
		Count(&count).Error
	return int(count), err
}

// --- Expenses ---

func (r *repo) CreateExpenseWithShares(expense *TravelExpense, shares []TravelExpenseShare) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if expense.ID == uuid.Nil {
			expense.ID = uuid.New()
		}
		if err := tx.Create(expense).Error; err != nil {
			return err
		}
		for i := range shares {
			if shares[i].ID == uuid.Nil {
				shares[i].ID = uuid.New()
			}
			shares[i].ExpenseID = expense.ID
			if err := tx.Create(&shares[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repo) GetExpense(id uuid.UUID) (*TravelExpense, error) {
	var e TravelExpense
	err := r.db.First(&e, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrExpenseNotFound
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repo) ListExpensesByGroup(groupID uuid.UUID) ([]TravelExpense, error) {
	var es []TravelExpense
	err := r.db.Where("group_id = ?", groupID).
		Order("date DESC, created_at DESC").
		Find(&es).Error
	return es, err
}

func (r *repo) DeleteExpense(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&TravelExpense{}).Error
}

func (r *repo) ListSharesByExpense(expenseID uuid.UUID) ([]TravelExpenseShare, error) {
	var ss []TravelExpenseShare
	err := r.db.Where("expense_id = ?", expenseID).Find(&ss).Error
	return ss, err
}

// --- Settlements ---

func (r *repo) CreateSettlement(s *TravelSettlement) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.Status == "" {
		s.Status = SettlementPending
	}
	return r.db.Create(s).Error
}

func (r *repo) GetSettlement(id uuid.UUID) (*TravelSettlement, error) {
	var s TravelSettlement
	err := r.db.First(&s, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrExpenseNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repo) ListSettlementsByGroup(groupID uuid.UUID) ([]TravelSettlement, error) {
	var ss []TravelSettlement
	err := r.db.Where("group_id = ?", groupID).
		Order("created_at DESC").
		Find(&ss).Error
	return ss, err
}

func (r *repo) UpdateSettlement(s *TravelSettlement) error {
	return r.db.Save(s).Error
}

// --- Aggregations ---

func (r *repo) SumPaidByUser(groupID, userID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.Model(&TravelExpense{}).
		Where("group_id = ? AND paid_by = ?", groupID, userID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *repo) SumShareByUser(groupID, userID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.Model(&TravelExpenseShare{}).
		Joins("JOIN travel_expenses e ON e.id = travel_expense_shares.expense_id").
		Where("e.group_id = ? AND travel_expense_shares.user_id = ?", groupID, userID).
		Select("COALESCE(SUM(travel_expense_shares.amount), 0)").
		Scan(&total).Error
	return total, err
}

// BalancesByGroup returns, in a single round-trip per side, the net
// balance of every member in a travel group. Replaces the N+1 pattern
// of calling SumPaidByUser + SumShareByUser inside a per-member loop.
type MemberBalance struct {
	UserID uuid.UUID
	Paid   int64
	Owed   int64
}

func (r *repo) BalancesByGroup(groupID uuid.UUID) ([]MemberBalance, error) {
	// One union query: sum(paid) per (group, paid_by) minus sum(share) per
	// (group, share.user_id). LEFT JOIN to members keeps rows even when
	// a user has never paid and never owed.
	rows, err := r.db.Raw(`
		WITH paid AS (
			SELECT paid_by AS user_id, SUM(amount) AS paid
			FROM travel_expenses
			WHERE group_id = ?
			GROUP BY paid_by
		), owed AS (
			SELECT s.user_id, SUM(s.amount) AS owed
			FROM travel_expense_shares s
			JOIN travel_expenses e ON e.id = s.expense_id
			WHERE e.group_id = ?
			GROUP BY s.user_id
		), members AS (
			SELECT user_id FROM travel_group_members WHERE group_id = ?
		)
		SELECT m.user_id,
		       COALESCE(p.paid, 0) AS paid,
		       COALESCE(o.owed, 0) AS owed,
		       COALESCE(p.paid, 0) - COALESCE(o.owed, 0) AS balance
		FROM members m
		LEFT JOIN paid p ON p.user_id = m.user_id
		LEFT JOIN owed o ON o.user_id = m.user_id`,
		groupID, groupID, groupID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]MemberBalance, 0)
	for rows.Next() {
		var b MemberBalance
		if err := rows.Scan(&b.UserID, &b.Paid, &b.Owed); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// ListExpensesWithSharesByGroup loads every expense for a group along
// with all its shares in two queries (vs the previous N+1 of
// ListSharesByExpense per expense). Use when rendering a full group
// view that needs both rows.
type ExpenseWithShares struct {
	Expense TravelExpense
	Shares  []TravelExpenseShare
}

func (r *repo) ListExpensesWithSharesByGroup(groupID uuid.UUID) ([]ExpenseWithShares, error) {
	expenses, err := r.ListExpensesByGroup(groupID)
	if err != nil {
		return nil, err
	}
	if len(expenses) == 0 {
		return nil, nil
	}
	ids := make([]uuid.UUID, len(expenses))
	expByID := make(map[uuid.UUID]*TravelExpense, len(expenses))
	for i := range expenses {
		ids[i] = expenses[i].ID
		expByID[expenses[i].ID] = &expenses[i]
	}
	var allShares []TravelExpenseShare
	if err := r.db.Where("expense_id IN ?", ids).Find(&allShares).Error; err != nil {
		return nil, err
	}
	sharesByExpense := make(map[uuid.UUID][]TravelExpenseShare, len(expenses))
	for _, sh := range allShares {
		sharesByExpense[sh.ExpenseID] = append(sharesByExpense[sh.ExpenseID], sh)
	}
	out := make([]ExpenseWithShares, 0, len(expenses))
	for i := range expenses {
		out = append(out, ExpenseWithShares{
			Expense: expenses[i],
			Shares:  sharesByExpense[expenses[i].ID],
		})
	}
	return out, nil
}

// isUniqueViolation checks whether a Postgres error is the UNIQUE
// constraint violation (SQLSTATE 23505). Exposed here so the import lives
// next to its only consumer.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	// GORM may return a wrapped/driver-agnostic error.
	return err != nil && (err.Error() == "ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)" ||
		err.Error() == "UNIQUE constraint failed: travel_group_members.group_id, travel_group_members.user_id")
}