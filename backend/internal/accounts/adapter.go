package accounts

import (
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