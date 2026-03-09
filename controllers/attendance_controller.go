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

type AttendanceController struct {
	service services.AttendanceService
}

func NewAttendanceController(service services.AttendanceService) *AttendanceController {
	return &AttendanceController{service: service}
}

type attendanceRequest struct {
	Identifier string `json:"identifier" binding:"required"`
}

func (ctl *AttendanceController) CheckIn(c *gin.Context) {
	var req attendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	authUser := middleware.MustGetAuthUser(c)
	record, err := ctl.service.CheckIn(req.Identifier, authUser.Role, authUser.CenterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": record, "message": "time in marked"})
}

func (ctl *AttendanceController) CheckOut(c *gin.Context) {
	var req attendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	authUser := middleware.MustGetAuthUser(c)
	record, err := ctl.service.CheckOut(req.Identifier, authUser.Role, authUser.CenterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": record, "message": "time out marked"})
}

func (ctl *AttendanceController) Scan(c *gin.Context) {
	var req attendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	authUser := middleware.MustGetAuthUser(c)
	record, action, err := ctl.service.Scan(req.Identifier, authUser.Role, authUser.CenterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": record, "action": action})
}

func (ctl *AttendanceController) List(c *gin.Context) {
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
	if value := c.Query("center_id"); value != "" {
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
	records, err := ctl.service.List(filter, authUser.Role, authUser.CenterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": records})
}

func (ctl *AttendanceController) Dashboard(c *gin.Context) {
	authUser := middleware.MustGetAuthUser(c)
	response, err := ctl.service.Dashboard(authUser.Role, authUser.CenterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": response})
}
