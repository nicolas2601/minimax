package transactions

import (
	"time"
)

type CreateRequestDTO struct {
	AccountID   string  `json:"account_id" binding:"required,uuid"`
	CategoryID  *string `json:"category_id,omitempty" binding:"omitempty,uuid"`
	Type        string  `json:"type" binding:"required,oneof=expense income"`
	Amount      int64   `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"omitempty,len=3"`
	Date        string  `json:"date"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=255"`
	Notes       *string `json:"notes,omitempty"`
}

type UpdateRequestDTO struct {
	CategoryID  *string `json:"category_id,omitempty" binding:"omitempty,uuid"`
	Amount      *int64  `json:"amount,omitempty" binding:"omitempty,gt=0"`
	Date        *string `json:"date"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=255"`
	Notes       *string `json:"notes,omitempty"`
}

type TransferRequestDTO struct {
	FromAccountID string  `json:"from_account_id" binding:"required,uuid"`
	ToAccountID   string  `json:"to_account_id" binding:"required,uuid"`
	Amount        int64   `json:"amount" binding:"required,gt=0"`
	Currency      string  `json:"currency" binding:"omitempty,len=3"`
	Date          string  `json:"date"`
	Description   *string `json:"description,omitempty" binding:"omitempty,max=255"`
	Notes         *string `json:"notes,omitempty"`
}

type ListFilterQuery struct {
	AccountID  string `form:"account_id" binding:"omitempty,uuid"`
	CategoryID string `form:"category_id" binding:"omitempty,uuid"`
	Type       string `form:"type" binding:"omitempty,oneof=expense income transfer"`
	From       string `form:"from"`
	To         string `form:"to"`
	Limit      int    `form:"limit" binding:"omitempty,min=1,max=500"`
}

type ListResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type TransferResponse struct {
	Transfer TransferResult `json:"transfer"`
}

// ToServiceCreate converts the wire DTO into the service-layer CreateRequest.
func (d CreateRequestDTO) ToServiceCreate() CreateRequest {
	return CreateRequest{
		AccountID:   d.AccountID,
		CategoryID:  d.CategoryID,
		Type:        d.Type,
		Amount:      d.Amount,
		Currency:    d.Currency,
		Date:        d.Date,
		Description: d.Description,
		Notes:       d.Notes,
	}
}

// ToServiceUpdate converts the wire DTO into the service-layer UpdateRequest.
func (d UpdateRequestDTO) ToServiceUpdate() UpdateRequest {
	return UpdateRequest{
		CategoryID:  d.CategoryID,
		Amount:      d.Amount,
		Date:        d.Date,
		Description: d.Description,
		Notes:       d.Notes,
	}
}

// ToServiceTransfer converts the wire DTO into the service-layer TransferRequest.
func (d TransferRequestDTO) ToServiceTransfer() TransferRequest {
	return TransferRequest{
		FromAccountID: d.FromAccountID,
		ToAccountID:   d.ToAccountID,
		Amount:        d.Amount,
		Currency:      d.Currency,
		Date:          d.Date,
		Description:   d.Description,
		Notes:         d.Notes,
	}
}

// ToListFilter parses query params into the service-layer ListFilter.
func (q ListFilterQuery) ToListFilter() ListFilter {
	f := ListFilter{Type: q.Type, Limit: q.Limit}
	if q.AccountID != "" {
		id := parseUUID(q.AccountID)
		f.AccountID = &id
	}
	if q.CategoryID != "" {
		id := parseUUID(q.CategoryID)
		f.CategoryID = &id
	}
	if q.From != "" {
		if t, err := time.Parse("2006-01-02", q.From); err == nil {
			f.From = &t
		}
	}
	if q.To != "" {
		if t, err := time.Parse("2006-01-02", q.To); err == nil {
			f.To = &t
		}
	}
	return f
}