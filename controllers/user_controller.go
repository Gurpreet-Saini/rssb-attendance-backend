package controllers

import (
	"net/http"

	"attendance-system/backend/middleware"
	"attendance-system/backend/models"
	"attendance-system/backend/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

func (ctl *UserController) List(c *gin.Context) {
	authUser := middleware.MustGetAuthUser(c)
	role := c.Query("role")
	users, err := ctl.service.List(authUser.Role, authUser.CenterID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (ctl *UserController) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required"`
		CenterID *uint  `json:"center_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.MustGetAuthUser(c)
	if authUser.Role == models.RoleCenterAdmin {
		req.CenterID = authUser.CenterID
		if req.Role == models.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			return
		}
	}

	user, err := ctl.service.Create(req.Name, req.Email, req.Password, req.Role, req.CenterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": user})
}
