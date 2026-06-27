package categories

import (
	"errors"

	"github.com/google/uuid"
)

// TransactionCategoryAdapter exposes the subset of categories.Service that
// the transactions package needs: validate that a category exists and is
// owned by the user before attaching it to a transaction.
//
// We keep the contracts as tiny private interfaces inside each adapter so
// the categories package does not have to import the consumers (which
// would cause import cycles).
type TransactionCategoryAdapter struct {
	svc Service
}

// categoryLookup is the contract implemented by this adapter. We declare it
// here so the constructor signature can mention it without pulling in the
// transactions package.
type categoryLookup interface {
	GetByID(id, userID uuid.UUID) (exists bool, err error)
}

// NewTransactionCategoryAdapter wires categories.Service into the
// transactions.CategoryLookup contract.
func NewTransactionCategoryAdapter(svc Service) *TransactionCategoryAdapter {
	if svc == nil {
		return nil
	}
	return &TransactionCategoryAdapter{svc: svc}
}

// GetByID satisfies the look-up contract. It returns false when the
// category does not exist (or belongs to a different user) so the service
// layer can map that to its own ErrCategoryNotFound.
func (a *TransactionCategoryAdapter) GetByID(id, userID uuid.UUID) (bool, error) {
	_, err := a.svc.Get(id, userID)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, ErrCategoryNotFound) {
		return false, nil
	}
	return false, err
}

// BudgetCategoryAdapter exposes the same GetByID contract for the budgets
// package.
type BudgetCategoryAdapter struct {
	svc Service
}

// NewBudgetCategoryAdapter wires categories.Service into the budgets
// CategoryLookup contract. Separate constructor so the wiring stays explicit.
func NewBudgetCategoryAdapter(svc Service) *BudgetCategoryAdapter {
	if svc == nil {
		return nil
	}
	return &BudgetCategoryAdapter{svc: svc}
}

// GetByID satisfies the budget look-up contract.
func (a *BudgetCategoryAdapter) GetByID(id, userID uuid.UUID) (bool, error) {
	_, err := a.svc.Get(id, userID)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, ErrCategoryNotFound) {
		return false, nil
	}
	return false, err
}

// compile-time assertions that both adapters satisfy the tiny look-up
// contract. The actual transactions.Account / budgets.CategoryLookup types
// live in the consumer packages — the assertions below make sure we keep
// the signatures in sync.
var (
	_ categoryLookup = (*TransactionCategoryAdapter)(nil)
	_ categoryLookup = (*BudgetCategoryAdapter)(nil)
)