package categories

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidType = errors.New("invalid category type")
)

type Service interface {
	Create(userID uuid.UUID, req CreateRequest) (*Category, error)
	List(userID uuid.UUID, categoryType string) ([]Category, error)
	Get(id, userID uuid.UUID) (*Category, error)
	Update(id, userID uuid.UUID, req UpdateRequest) (*Category, error)
	Delete(id, userID uuid.UUID) error
	SeedDefaults(userID uuid.UUID) (int, error) // returns number of categories created
}

type service struct {
	repo CategoryRepository
}

func NewService(repo CategoryRepository) Service {
	return &service{repo: repo}
}

func (s *service) Create(userID uuid.UUID, req CreateRequest) (*Category, error) {
	if !IsValidType(req.Type) {
		return nil, ErrInvalidType
	}
	c := &Category{
		UserID: userID,
		Name:   req.Name,
		Type:   CategoryType(req.Type),
		Icon:   req.Icon,
		Color:  req.Color,
	}
	if req.ParentID != nil {
		pid, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, err
		}
		c.ParentID = &pid
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) List(userID uuid.UUID, categoryType string) ([]Category, error) {
	return s.repo.ListByUser(userID, categoryType)
}

func (s *service) Get(id, userID uuid.UUID) (*Category, error) {
	return s.repo.GetByID(id, userID)
}

func (s *service) Update(id, userID uuid.UUID, req UpdateRequest) (*Category, error) {
	c, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Color != nil {
		c.Color = req.Color
	}
	if req.Icon != nil {
		c.Icon = req.Icon
	}
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) Delete(id, userID uuid.UUID) error {
	return s.repo.Delete(id, userID)
}

func (s *service) SeedDefaults(userID uuid.UUID) (int, error) {
	created := 0
	for _, def := range DefaultCategoriesES {
		exists, err := s.repo.ExistsByNameAndType(userID, def.Name, string(def.Type))
		if err != nil {
			return created, err
		}
		if exists {
			continue
		}
		c := def
		c.UserID = userID
		c.IsDefault = true
		if err := s.repo.Create(&c); err != nil {
			return created, err
		}
		created++
	}
	return created, nil
}

func IsValidType(t string) bool {
	switch CategoryType(t) {
	case TypeExpense, TypeIncome:
		return true
	}
	return false
}