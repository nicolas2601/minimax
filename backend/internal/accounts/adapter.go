package accounts

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nicolas/finanzas/backend/internal/transactions"
)

// TransactionAccountAdapter exposes the subset of accounts.Service that the
// transactions package needs (validate ownership and currency when creating
// or transferring transactions). Keeping this in the accounts package — and
// not in transactions — avoids transactions importing accounts directly,
// which would let accounts drift toward a dependency on transactions later.
type TransactionAccountAdapter struct {
	svc Service
}

// NewTransactionAccountAdapter wires accounts.Service into the
// transactions.Account contract. Returns nil if svc is nil so main.go can
// pass nil to disable ownership checks (e.g. in tests).
func NewTransactionAccountAdapter(svc Service) *TransactionAccountAdapter {
	if svc == nil {
		return nil
	}
	return &TransactionAccountAdapter{svc: svc}
}

// GetByID implements transactions.Account. We map the accounts.Account
// into the transactions.AccountInfo DTO so the consumer stays decoupled.
func (a *TransactionAccountAdapter) GetByID(id, userID uuid.UUID) (transactions.AccountInfo, error) {
	acc, err := a.svc.Get(id, userID)
	if err != nil {
		return transactions.AccountInfo{}, err
	}
	return transactions.AccountInfo{
		ID:       acc.ID,
		UserID:   acc.UserID,
		Currency: acc.Currency,
	}, nil
}

// RecurringAccountAdapter exposes the same Exists contract for the
// recurring package: does this account exist and is it owned by the user.
// Same shape as GoalsAccountAdapter so they could be unified later, but kept
// as separate types today because their consumers (recurring.AccountLookup,
// goals.AccountLookup) live in different packages and a single adapter would
// reintroduce the import cycle we're trying to avoid.
type RecurringAccountAdapter struct {
	svc Service
}

// NewRecurringAccountAdapter wires accounts.Service into the
// recurring.AccountLookup contract.
func NewRecurringAccountAdapter(svc Service) *RecurringAccountAdapter {
	if svc == nil {
		return nil
	}
	return &RecurringAccountAdapter{svc: svc}
}

// Exists implements recurring.AccountLookup. A "not found" result maps to
// (false, nil) so the recurring service can distinguish missing from real DB
// errors, which propagate.
func (a *RecurringAccountAdapter) Exists(id, userID uuid.UUID) (bool, error) {
	_, err := a.svc.Get(id, userID)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}