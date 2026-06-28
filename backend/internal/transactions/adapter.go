package transactions

import (
	"time"

	"github.com/google/uuid"
)

// RecurringTxCreatorAdapter exposes the subset of transactions.Service that
// the recurring package needs to materialize one transaction per scheduled
// occurrence. Keeping the adapter in the transactions package avoids the
// recurring package importing transactions directly (which would create a
// dependency arrow pointing the wrong way for clean architecture).
type RecurringTxCreatorAdapter struct {
	svc *Service
}

// NewRecurringTxCreatorAdapter wires transactions.Service into the
// recurring.TxCreator contract. Returns nil if svc is nil so main.go can
// pass nil to disable the recurring engine (e.g. in tests).
func NewRecurringTxCreatorAdapter(svc *Service) *RecurringTxCreatorAdapter {
	if svc == nil {
		return nil
	}
	return &RecurringTxCreatorAdapter{svc: svc}
}

// CreateFromRecurring implements recurring.TxCreator. It delegates directly
// to the service method of the same name — the adapter exists purely to
// keep the import graph one-way (transactions never imports recurring).
func (a *RecurringTxCreatorAdapter) CreateFromRecurring(
	userID, accountID, categoryID uuid.UUID,
	txType string,
	amount int64,
	currency string,
	date time.Time,
	description, notes *string,
	recurringRunID uuid.UUID,
) (uuid.UUID, error) {
	return a.svc.CreateFromRecurring(
		userID, accountID, categoryID,
		txType, amount, currency,
		date, description, notes,
		recurringRunID,
	)
}