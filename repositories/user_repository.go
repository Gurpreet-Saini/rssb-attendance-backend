package repositories

import (
	"attendance-system/backend/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByResetToken(token string) (*models.User, error)
	List(centerID *uint, role string) ([]models.User, error)
	Create(user *models.User) error
	Save(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Center").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Center").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByResetToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Center").Where("password_reset_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) List(centerID *uint, role string) ([]models.User, error) {
	var users []models.User
	query := r.db.Preload("Center").Order("created_at desc")
	if centerID != nil {
		query = query.Where("center_id = ?", *centerID)
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	err := query.Find(&users).Error
	return users, err
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Save(user *models.User) error {
	return r.db.Save(user).Error
}
