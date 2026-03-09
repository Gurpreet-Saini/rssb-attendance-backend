package controllers

import (
	"net/http"
	"strconv"
	"time"

	"attendance-system/backend/middleware"
	"attendance-system/backend/repositories"
	"attendance-system/backend/services"
	"github.com/gin-gonic/gin"
)

type ReportController struct {
	service services.ReportService
}

func NewReportController(service services.ReportService) *ReportController {
	return &ReportController{service: service}
}

func (ctl *ReportController) ExportExcel(c *gin.Context) {
	authUser := middleware.MustGetAuthUser(c)
	filter := repositories.AttendanceFilter{}
	if value := c.Query("start_date"); value != "" {
		if parsed, err := time.Parse("2006-01-02", value); err == nil {
			filter.StartDate = &parsed
		}
	}
	if value := c.Query("end_date"); value != "" {
		if parsed, err := time.Parse("2006-01-02", value); err == nil {
			filter.EndDate = &parsed
		}
	}
	if value := c.Query("center_id"); value != "" && authUser.Role == "super_admin" {
		if parsed, err := strconv.ParseUint(value, 10, 64); err == nil {
			centerID := uint(parsed)
			filter.CenterID = &centerID
		}
	}
	if value := c.Query("employee_id"); value != "" {
		if parsed, err := strconv.ParseUint(value, 10, 64); err == nil {
			employeeID := uint(parsed)
			filter.EmployeeID = &employeeID
		}
	}
	file, err := ctl.service.GenerateAttendanceExcel(filter, authUser.Role, authUser.CenterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", `attachment; filename="attendance-report.xlsx"`)
	c.Header("File-Name", "attendance-report.xlsx")
	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
}
