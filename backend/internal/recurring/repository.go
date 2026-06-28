package recurring

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRuleNotFound = errors.New("recurring rule not found")
	ErrRunNotFound  = errors.New("recurring run not found")
)

type Repository interface {
	// Rules
	CreateRule(r *Rule) error
	GetRuleByID(id, userID uuid.UUID) (*Rule, error)
	ListRulesByUser(userID uuid.UUID) ([]Rule, error)
	UpdateRule(r *Rule) error
	DeleteRule(id, userID uuid.UUID) error
	ListDueRules(before time.Time) ([]Rule, error) // is_active=true AND next_run_date <= before

	// Runs
	CreateRun(run *Run) error
	GetRunByRuleAndDate(ruleID uuid.UUID, scheduled time.Time) (*Run, error)
	ListRunsByRule(ruleID uuid.UUID, limit int) ([]Run, error)
	UpdateRun(run *Run) error
}

type repo struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) CreateRule(rr *Rule) error {
	if rr.ID == uuid.Nil {
		rr.ID = uuid.New()
	}
	return r.db.Create(rr).Error
}

func (r *repo) GetRuleByID(id, userID uuid.UUID) (*Rule, error) {
	var rr Rule
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&rr).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRuleNotFound
	}
	if err != nil {
		return nil, err
	}
	return &rr, nil
}

func (r *repo) ListRulesByUser(userID uuid.UUID) ([]Rule, error) {
	var list []Rule
	err := r.db.Where("user_id = ?", userID).Order("next_run_date ASC").Find(&list).Error
	return list, err
}

func (r *repo) UpdateRule(rr *Rule) error {
	return r.db.Save(rr).Error
}

func (r *repo) DeleteRule(id, userID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&Rule{}).Error
}

func (r *repo) ListDueRules(before time.Time) ([]Rule, error) {
	var list []Rule
	err := r.db.Where("is_active = ? AND next_run_date <= ?", true, before).
		Order("next_run_date ASC").
		Find(&list).Error
	return list, err
}

func (r *repo) CreateRun(run *Run) error {
	if run.ID == uuid.Nil {
		run.ID = uuid.New()
	}
	return r.db.Create(run).Error
}

func (r *repo) GetRunByRuleAndDate(ruleID uuid.UUID, scheduled time.Time) (*Run, error) {
	var run Run
	err := r.db.Where("recurring_rule_id = ? AND scheduled_date = ?", ruleID, scheduled).
		First(&run).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRunNotFound
	}
	if err != nil {
		return nil, err
	}
	return &run, nil
}

func (r *repo) ListRunsByRule(ruleID uuid.UUID, limit int) ([]Run, error) {
	var list []Run
	q := r.db.Where("recurring_rule_id = ?", ruleID).Order("scheduled_date DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *repo) UpdateRun(run *Run) error {
	return r.db.Save(run).Error
}
