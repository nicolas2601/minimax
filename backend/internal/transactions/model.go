package transactions

import (
	"time"

	"github.com/google/uuid"
)

// TxType represents the kind of transaction. Transfers are represented as
// TWO transactions linked by TransferPairID — one for the source account
// (expense-style) and one for the destination account (income-style).
type TxType string

const (
	TypeExpense  TxType = "expense"
	TypeIncome   TxType = "income"
	TypeTransfer TxType = "transfer"
)

type Transaction struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	AccountID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"account_id"`
	CategoryID     *uuid.UUID `gorm:"type:uuid;index" json:"category_id,omitempty"`
	Type           TxType     `gorm:"not null;size:20" json:"type"`
	Amount         int64      `gorm:"not null" json:"amount"`
	Currency       string     `gorm:"not null;size:3;default:COP" json:"currency"`
	Date           time.Time  `gorm:"not null;index" json:"date"`
	Description    *string    `gorm:"size:255" json:"description,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	TransferPairID *uuid.UUID `gorm:"type:uuid;index" json:"transfer_pair_id,omitempty"`
	DeletedAt      *time.Time `gorm:"index" json:"-"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}

func IsValidType(t string) bool {
	switch TxType(t) {
	case TypeExpense, TypeIncome, TypeTransfer:
		return true
	}
	return false
}