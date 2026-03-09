package services

import (
	"attendance-system/backend/repositories"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type ReportService interface {
	GenerateAttendanceExcel(recordsFilter repositories.AttendanceFilter, requesterRole string, requesterCenterID *uint) (*excelize.File, error)
}

type reportService struct {
	attendance AttendanceService
}

func NewReportService(attendance AttendanceService) ReportService {
	return &reportService{attendance: attendance}
}

func (s *reportService) GenerateAttendanceExcel(recordsFilter repositories.AttendanceFilter, requesterRole string, requesterCenterID *uint) (*excelize.File, error) {
	records, err := s.attendance.List(recordsFilter, requesterRole, requesterCenterID)
	if err != nil {
		return nil, err
	}

	file := excelize.NewFile()
	sheet := "Attendance"
	file.SetSheetName("Sheet1", sheet)

	headers := []string{"Date", "Employee Name", "Employee ID", "Badge Number", "Center", "Time In", "Time Out", "Status"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheet, cell, header)
	}

	for i, record := range records {
		row := i + 2
		timeIn := ""
		if record.TimeIn != nil {
			timeIn = record.TimeIn.Format(time.RFC3339)
		}
		timeOut := ""
		if record.TimeOut != nil {
			timeOut = record.TimeOut.Format(time.RFC3339)
		}
		values := []interface{}{
			record.Date.Format("2006-01-02"),
			record.Employee.Name,
			record.Employee.EmployeeID,
			record.Employee.BadgeNumber,
			record.Center.Name,
			timeIn,
			timeOut,
			record.Status,
		}
		for col, value := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row)
			file.SetCellValue(sheet, cell, value)
		}
	}

	for col := 1; col <= len(headers); col++ {
		column, _ := excelize.ColumnNumberToName(col)
		file.SetColWidth(sheet, column, column, 20)
	}
	file.SetDocProps(&excelize.DocProperties{
		Creator:     "Attendance Management System",
		Description: "Attendance export",
		Identifier:  fmt.Sprintf("attendance-report-%d", time.Now().Unix()),
		Title:       "Attendance Report",
	})

	return file, nil
}
