package repositories

import (
	"attendance-system/backend/models"
	"gorm.io/gorm"
)

type CenterRepository interface {
	List() ([]models.Center, error)
	Create(center *models.Center) error
	GetByID(id uint) (*models.Center, error)
	SoftDeleteCascade(id uint) error
}

type centerRepository struct {
	db *gorm.DB
}

func NewCenterRepository(db *gorm.DB) CenterRepository {
	return &centerRepository{db: db}
}

func (r *centerRepository) List() ([]models.Center, error) {
	var centers []models.Center
	err := r.db.Order("name asc").Find(&centers).Error
	return centers, err
}

func (r *centerRepository) Create(center *models.Center) error {
	return r.db.Create(center).Error
}

func (r *centerRepository) GetByID(id uint) (*models.Center, error) {
	var center models.Center
	err := r.db.First(&center, id).Error
	if err != nil {
		return nil, err
	}
	return &center, nil
}

func (r *centerRepository) SoftDeleteCascade(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("center_id = ?", id).Delete(&models.Attendance{}).Error; err != nil {
			return err
		}
		if err := tx.Where("center_id = ?", id).Delete(&models.Employee{}).Error; err != nil {
			return err
		}
		if err := tx.Where("center_id = ?", id).Delete(&models.User{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Center{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
