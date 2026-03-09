package controllers

import (
	"net/http"
	"strconv"

	"attendance-system/backend/services"
	"github.com/gin-gonic/gin"
)

type CenterController struct {
	service services.CenterService
}

func NewCenterController(service services.CenterService) *CenterController {
	return &CenterController{service: service}
}

func (ctl *CenterController) List(c *gin.Context) {
	centers, err := ctl.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": centers})
}

func (ctl *CenterController) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Location string `json:"location" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	center, err := ctl.service.Create(req.Name, req.Location)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": center})
}

func (ctl *CenterController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}
	if err := ctl.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "center soft deleted with related users, employees, and attendance"})
}
