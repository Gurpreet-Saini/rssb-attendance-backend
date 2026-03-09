package controllers

import (
	"net/http"
	"strconv"

	"attendance-system/backend/middleware"
	"attendance-system/backend/models"
	"attendance-system/backend/services"
	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	service services.EmployeeService
}

func NewEmployeeController(service services.EmployeeService) *EmployeeController {
	return &EmployeeController{service: service}
}

func (ctl *EmployeeController) List(c *gin.Context) {
	authUser := middleware.MustGetAuthUser(c)
	var centerID *uint
	if raw := c.Query("center_id"); raw != "" {
		if parsed, err := strconv.ParseUint(raw, 10, 64); err == nil {
			value := uint(parsed)
			centerID = &value
		}
	}
	employees, err := ctl.service.List(c.Query("q"), authUser.Role, authUser.CenterID, centerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": employees})
}

func (ctl *EmployeeController) Create(c *gin.Context) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.MustGetAuthUser(c)
	if authUser.Role == models.RoleCenterAdmin {
		employee.CenterID = *authUser.CenterID
	}
	if authUser.Role == models.RoleOperator {
		c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}

	if err := ctl.service.Create(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": employee})
}

func (ctl *EmployeeController) BulkUpload(c *gin.Context) {
	authUser := middleware.MustGetAuthUser(c)
	centerIDValue, err := strconv.ParseUint(c.PostForm("center_id"), 10, 64)
	if err != nil && authUser.Role == models.RoleSuperAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"message": "center_id is required"})
		return
	}

	var centerID uint
	if authUser.Role == models.RoleCenterAdmin {
		centerID = *authUser.CenterID
	} else {
		centerID = uint(centerIDValue)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "file is required"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	defer file.Close()

	count, err := ctl.service.BulkCreate(file, centerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "employees uploaded", "count": count})
}

func (ctl *EmployeeController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}
	var payload models.Employee
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	authUser := middleware.MustGetAuthUser(c)
	employee, err := ctl.service.Update(uint(id), &payload, authUser.Role, authUser.CenterID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": employee})
}

func (ctl *EmployeeController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}
	authUser := middleware.MustGetAuthUser(c)
	if err := ctl.service.Delete(uint(id), authUser.Role, authUser.CenterID); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "employee deleted"})
}
