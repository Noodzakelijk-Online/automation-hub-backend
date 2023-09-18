package automation

import (
	"automation-hub-backend/internal/infra"
	"automation-hub-backend/internal/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindByID(id uuid.UUID) (*model.Automation, error)
	Create(automation *model.Automation) (*model.Automation, error)
	Update(automation *model.Automation) (*model.Automation, error)
	Delete(id uuid.UUID) error
	FindAll() ([]*model.Automation, error)
	MaxPosition() (int, error)
	Transaction(txFunc func(tx *gorm.DB) error) (err error)
}

type GormUserRepository struct {
	DB *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) Repository {
	return &GormUserRepository{
		DB: db,
	}
}

func DefaultRepository() Repository {
	db, err := infra.GetDefaultDB()
	if err != nil {
		panic(err)
	}
	return NewGormUserRepository(db)
}

func (r *GormUserRepository) FindByID(id uuid.UUID) (*model.Automation, error) {
	var automation model.Automation
	err := r.DB.First(&automation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &automation, nil
}

func (r *GormUserRepository) Create(automation *model.Automation) (*model.Automation, error) {
	err := r.DB.Create(automation).Error
	if err != nil {
		return nil, err
	}
	return automation, nil
}

func (r *GormUserRepository) Update(automation *model.Automation) (*model.Automation, error) {
	err := r.DB.Save(automation).Error
	if err != nil {
		return nil, err
	}
	return automation, nil
}

func (r *GormUserRepository) Delete(id uuid.UUID) error {
	err := r.DB.Delete(&model.Automation{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormUserRepository) FindAll() ([]*model.Automation, error) {
	var automations []*model.Automation
	err := r.DB.Find(&automations).Error
	if err != nil {
		return nil, err
	}
	return automations, nil
}

func (r *GormUserRepository) Transaction(txFunc func(tx *gorm.DB) error) (err error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("transaction panicked: %v", r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	err = txFunc(tx)
	return err
}

func (r *GormUserRepository) MaxPosition() (int, error) {
	var automation model.Automation
	err := r.DB.Order("position desc").First(&automation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return automation.Position, nil
}
