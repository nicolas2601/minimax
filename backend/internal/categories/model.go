package categories

import (
	"time"

	"github.com/google/uuid"
)

type CategoryType string

const (
	TypeExpense CategoryType = "expense"
	TypeIncome  CategoryType = "income"
)

type Category struct {
	ID        uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string       `gorm:"not null;size:100" json:"name"`
	Type      CategoryType `gorm:"not null;size:20" json:"type"`
	ParentID  *uuid.UUID   `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	Icon      *string      `gorm:"size:50" json:"icon,omitempty"`
	Color     *string      `gorm:"size:7" json:"color,omitempty"`
	IsDefault bool         `gorm:"not null;default:false" json:"is_default"`
	DeletedAt *time.Time   `gorm:"index" json:"-"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}