package services

import (
	"attendance-system/backend/models"
	"attendance-system/backend/repositories"
	"fmt"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"
)

type EmployeeService interface {
	List(query string, requesterRole string, requesterCenterID *uint, filterCenterID *uint) ([]models.Employee, error)
	Create(employee *models.Employee) error
	BulkCreate(reader io.Reader, centerID uint) (int, error)
	Update(id uint, payload *models.Employee, requesterRole string, requesterCenterID *uint) (*models.Employee, error)
	Delete(id uint, requesterRole string, requesterCenterID *uint) error
	FindByIdentifier(identifier string, centerID *uint) (*models.Employee, error)
}

type employeeService struct {
	repo   repositories.EmployeeRepository
	audits repositories.AuditRepository
}

func NewEmployeeService(repo repositories.EmployeeRepository, audits repositories.AuditRepository) EmployeeService {
	return &employeeService{repo: repo, audits: audits}
}

func (s *employeeService) List(query string, requesterRole string, requesterCenterID *uint, filterCenterID *uint) ([]models.Employee, error) {
	centerID := filterCenterID
	if requesterRole == models.RoleCenterAdmin || requesterRole == models.RoleOperator {
		centerID = requesterCenterID
	}
	return s.repo.List(repositories.EmployeeFilter{Query: query, CenterID: centerID})
}

func (s *employeeService) Create(employee *models.Employee) error {
	if err := s.repo.Create(employee); err != nil {
		return err
	}
	return s.audits.Create(&models.AuditLog{
		Action:     "employee_created",
		EntityType: "employee",
		EntityID:   fmt.Sprintf("%d", employee.ID),
		Details:    employee.EmployeeID,
	})
}

func (s *employeeService) BulkCreate(reader io.Reader, centerID uint) (int, error) {
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = file.Close()
	}()

	sheet := file.GetSheetName(0)
	if sheet == "" {
		return 0, fmt.Errorf("excel sheet is empty")
	}

	rows, err := file.GetRows(sheet)
	if err != nil {
		return 0, err
	}

	employees := make([]models.Employee, 0, len(rows))
	for index, row := range rows {
		if len(row) == 0 {
			continue
		}
		if index == 0 && strings.EqualFold(strings.TrimSpace(row[0]), "name") {
			continue
		}
		if len(row) < 5 {
			return 0, fmt.Errorf("row %d must contain name, employee_id, badge_number, barcode, designation", index+1)
		}
		employees = append(employees, models.Employee{
			Name:        strings.TrimSpace(row[0]),
			EmployeeID:  strings.TrimSpace(row[1]),
			BadgeNumber: strings.TrimSpace(row[2]),
			Barcode:     strings.TrimSpace(row[3]),
			Designation: strings.TrimSpace(row[4]),
			CenterID:    centerID,
		})
	}

	if len(employees) == 0 {
		return 0, fmt.Errorf("no employee rows found")
	}
	if err := s.repo.BulkCreate(employees); err != nil {
		return 0, err
	}
	return len(employees), nil
}

func (s *employeeService) Update(id uint, payload *models.Employee, requesterRole string, requesterCenterID *uint) (*models.Employee, error) {
	employee, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if (requesterRole == models.RoleCenterAdmin || requesterRole == models.RoleOperator) && (requesterCenterID == nil || employee.CenterID != *requesterCenterID) {
		return nil, fmt.Errorf("forbidden")
	}
	employee.Name = payload.Name
	employee.EmployeeID = payload.EmployeeID
	employee.BadgeNumber = payload.BadgeNumber
	employee.Barcode = payload.Barcode
	employee.Designation = payload.Designation
	if requesterRole == models.RoleSuperAdmin {
		employee.CenterID = payload.CenterID
	}
	if err := s.repo.Update(employee); err != nil {
		return nil, err
	}
	return employee, nil
}

func (s *employeeService) Delete(id uint, requesterRole string, requesterCenterID *uint) error {
	employee, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if (requesterRole == models.RoleCenterAdmin || requesterRole == models.RoleOperator) && (requesterCenterID == nil || employee.CenterID != *requesterCenterID) {
		return fmt.Errorf("forbidden")
	}
	return s.repo.Delete(id)
}

func (s *employeeService) FindByIdentifier(identifier string, centerID *uint) (*models.Employee, error) {
	return s.repo.GetByIdentifier(identifier, centerID)
}
