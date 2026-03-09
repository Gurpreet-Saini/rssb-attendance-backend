package repositories

import (
	"attendance-system/backend/models"
	"strings"

	"gorm.io/gorm"
)

type EmployeeFilter struct {
	Query    string
	CenterID *uint
}

type EmployeeRepository interface {
	List(filter EmployeeFilter) ([]models.Employee, error)
	GetByID(id uint) (*models.Employee, error)
	GetByIdentifier(identifier string, centerID *uint) (*models.Employee, error)
	Create(employee *models.Employee) error
	BulkCreate(employees []models.Employee) error
	Update(employee *models.Employee) error
	Delete(id uint) error
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) List(filter EmployeeFilter) ([]models.Employee, error) {
	var employees []models.Employee
	query := r.db.Preload("Center").Order("created_at desc")
	if filter.CenterID != nil {
		query = query.Where("center_id = ?", *filter.CenterID)
	}
	if filter.Query != "" {
		pattern := "%" + strings.ToLower(filter.Query) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(employee_id) LIKE ? OR LOWER(badge_number) LIKE ? OR LOWER(barcode) LIKE ?",
			pattern, pattern, pattern, pattern,
		)
	}
	err := query.Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) GetByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.Preload("Center").First(&employee, id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) GetByIdentifier(identifier string, centerID *uint) (*models.Employee, error) {
	var employee models.Employee
	query := r.db.Preload("Center").Where(
		"employee_id = ? OR badge_number = ? OR barcode = ? OR LOWER(name) = LOWER(?)",
		identifier, identifier, identifier, identifier,
	)
	if centerID != nil {
		query = query.Where("center_id = ?", *centerID)
	}
	err := query.First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) Create(employee *models.Employee) error {
	return r.db.Create(employee).Error
}

func (r *employeeRepository) BulkCreate(employees []models.Employee) error {
	return r.db.Create(&employees).Error
}

func (r *employeeRepository) Update(employee *models.Employee) error {
	return r.db.Save(employee).Error
}

func (r *employeeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Employee{}, id).Error
}
