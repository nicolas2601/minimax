package categories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrCategoryNotFound = errors.New("category not found")

type CategoryRepository interface {
	Create(c *Category) error
	ListByUser(userID uuid.UUID, categoryType string) ([]Category, error)
	GetByID(id, userID uuid.UUID) (*Category, error)
	Update(c *Category) error
	Delete(id, userID uuid.UUID) error
	ExistsByNameAndType(userID uuid.UUID, name, categoryType string) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(c *Category) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return r.db.Create(c).Error
}

func (r *categoryRepository) ListByUser(userID uuid.UUID, categoryType string) ([]Category, error) {
	var list []Category
	q := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).Order("name ASC")
	if categoryType != "" {
		q = q.Where("type = ?", categoryType)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *categoryRepository) GetByID(id, userID uuid.UUID) (*Category, error) {
	var c Category
	err := r.db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).
		First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCategoryNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *categoryRepository) Update(c *Category) error {
	return r.db.Save(c).Error
}

func (r *categoryRepository) Delete(id, userID uuid.UUID) error {
	return r.db.Model(&Category{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *categoryRepository) ExistsByNameAndType(userID uuid.UUID, name, categoryType string) (bool, error) {
	var c Category
	err := r.db.Where("user_id = ? AND name = ? AND type = ? AND deleted_at IS NULL", userID, name, categoryType).
		First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}