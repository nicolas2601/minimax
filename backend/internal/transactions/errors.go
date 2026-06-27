package transactions

import "errors"

// Domain-level errors. Both the repository and the service return these so
// the handler can map them to HTTP status codes consistently.
var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrAccountMismatch     = errors.New("source and destination accounts must differ")
	ErrCurrencyMismatch    = errors.New("currency mismatch between accounts")
)