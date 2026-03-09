package services

import (
	"attendance-system/backend/models"
	"attendance-system/backend/repositories"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DashboardResponse struct {
	Counts      repositories.DashboardCounts `json:"counts"`
	CenterStats []repositories.CenterSummary `json:"center_stats"`
	DailyTrend  []repositories.DailyTrend    `json:"daily_trend"`
}

type AttendanceService interface {
	CheckIn(identifier string, requesterRole string, requesterCenterID *uint) (*models.Attendance, error)
	CheckOut(identifier string, requesterRole string, requesterCenterID *uint) (*models.Attendance, error)
	Scan(identifier string, requesterRole string, requesterCenterID *uint) (*models.Attendance, string, error)
	List(filter repositories.AttendanceFilter, requesterRole string, requesterCenterID *uint) ([]models.Attendance, error)
	Dashboard(requesterRole string, requesterCenterID *uint) (*DashboardResponse, error)
}

type attendanceService struct {
	attendance repositories.AttendanceRepository
	employees  repositories.EmployeeRepository
}

func NewAttendanceService(attendance repositories.AttendanceRepository, employees repositories.EmployeeRepository) AttendanceService {
	return &attendanceService{attendance: attendance, employees: employees}
}

func normalizeToday() time.Time {
	now := time.Now().In(time.FixedZone("IST", 5*3600+1800))
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func (s *attendanceService) CheckIn(identifier string, requesterRole string, requesterCenterID *uint) (*models.Attendance, error) {
	employee, err := s.resolveEmployee(identifier, requesterRole, requesterCenterID)
	if err != nil {
		return nil, err
	}
	today := normalizeToday()
	record, err := s.attendance.GetByEmployeeAndDate(employee.ID, today)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if record != nil && record.TimeIn != nil {
		return nil, fmt.Errorf("time in already marked")
	}

	now := time.Now().In(today.Location())
	if record == nil {
		record = &models.Attendance{
			EmployeeID: employee.ID,
			Date:       today,
			TimeIn:     &now,
			Status:     models.AttendancePartial,
			CenterID:   employee.CenterID,
		}
		if err := s.attendance.Create(record); err != nil {
			return nil, err
		}
		return record, nil
	}

	record.TimeIn = &now
	record.Status = deriveStatus(record.TimeIn, record.TimeOut)
	if err := s.attendance.Update(record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *attendanceService) CheckOut(identifier string, requesterRole string, requesterCenterID *uint) (*models.Attendance, error) {
	employee, err := s.resolveEmployee(identifier, requesterRole, requesterCenterID)
	if err != nil {
		return nil, err
	}
	today := normalizeToday()
	record, err := s.attendance.GetByEmployeeAndDate(employee.ID, today)
	if err != nil {
		return nil, fmt.Errorf("time in not marked")
	}
	if record.TimeOut != nil {
		return nil, fmt.Errorf("time out already marked")
	}
	now := time.Now().In(today.Location())
	record.TimeOut = &now
	record.Status = deriveStatus(record.TimeIn, record.TimeOut)
	if err := s.attendance.Update(record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *attendanceService) Scan(identifier string, requesterRole string, requesterCenterID *uint) (*models.Attendance, string, error) {
	employee, err := s.resolveEmployee(identifier, requesterRole, requesterCenterID)
	if err != nil {
		return nil, "", err
	}
	today := normalizeToday()
	record, err := s.attendance.GetByEmployeeAndDate(employee.ID, today)
	if errors.Is(err, gorm.ErrRecordNotFound) || record == nil || record.TimeIn == nil {
		result, err := s.CheckIn(identifier, requesterRole, requesterCenterID)
		return result, "checkin", err
	}
	result, err := s.CheckOut(identifier, requesterRole, requesterCenterID)
	return result, "checkout", err
}

func (s *attendanceService) List(filter repositories.AttendanceFilter, requesterRole string, requesterCenterID *uint) ([]models.Attendance, error) {
	if requesterRole == models.RoleCenterAdmin || requesterRole == models.RoleOperator {
		filter.CenterID = requesterCenterID
	}
	return s.attendance.List(filter)
}

func (s *attendanceService) Dashboard(requesterRole string, requesterCenterID *uint) (*DashboardResponse, error) {
	today := normalizeToday()
	lateAfter := time.Date(today.Year(), today.Month(), today.Day(), 9, 15, 0, 0, today.Location())
	centerID := requesterCenterID
	if requesterRole == models.RoleSuperAdmin {
		centerID = nil
	}

	counts, err := s.attendance.GetDashboardCounts(centerID, today, lateAfter)
	if err != nil {
		return nil, err
	}
	centerStats, err := s.attendance.GetCenterSummary(today)
	if err != nil {
		return nil, err
	}
	dailyTrend, err := s.attendance.GetDailyTrend(centerID, today.AddDate(0, 0, -6), today)
	if err != nil {
		return nil, err
	}
	return &DashboardResponse{
		Counts:      counts,
		CenterStats: centerStats,
		DailyTrend:  dailyTrend,
	}, nil
}

func (s *attendanceService) resolveEmployee(identifier string, requesterRole string, requesterCenterID *uint) (*models.Employee, error) {
	centerID := requesterCenterID
	if requesterRole == models.RoleSuperAdmin {
		centerID = nil
	}
	return s.employees.GetByIdentifier(identifier, centerID)
}

func deriveStatus(timeIn, timeOut *time.Time) string {
	if timeIn != nil && timeOut != nil {
		return models.AttendancePresent
	}
	if timeIn != nil || timeOut != nil {
		return models.AttendancePartial
	}
	return models.AttendanceAbsent
}
