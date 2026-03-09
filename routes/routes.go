package routes

import (
	"attendance-system/backend/controllers"
	"attendance-system/backend/middleware"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	Auth       *controllers.AuthController
	Centers    *controllers.CenterController
	Users      *controllers.UserController
	Employees  *controllers.EmployeeController
	Attendance *controllers.AttendanceController
	Reports    *controllers.ReportController
}

func Register(router *gin.Engine, secret string, ctl Controllers) {
	api := router.Group("/api")
	api.POST("/auth/login", ctl.Auth.Login)
	api.POST("/auth/forgot-password", ctl.Auth.ForgotPassword)
	api.POST("/auth/reset-password", ctl.Auth.ResetPassword)

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(secret))
	{
		protected.GET("/dashboard/summary", ctl.Attendance.Dashboard)
		protected.POST("/auth/change-password", ctl.Auth.ChangePassword)

		protected.GET("/centers", ctl.Centers.List)
		protected.POST("/centers", middleware.RequireRoles("super_admin"), ctl.Centers.Create)
		protected.DELETE("/centers/:id", middleware.RequireRoles("super_admin"), ctl.Centers.Delete)

		protected.GET("/users", middleware.RequireRoles("super_admin"), ctl.Users.List)
		protected.POST("/users", middleware.RequireRoles("super_admin"), ctl.Users.Create)

		protected.GET("/employees", middleware.RequireRoles("super_admin", "center_admin", "operator"), ctl.Employees.List)
		protected.POST("/employees", middleware.RequireRoles("super_admin", "center_admin"), ctl.Employees.Create)
		protected.POST("/employees/bulk-upload", middleware.RequireRoles("super_admin", "center_admin"), ctl.Employees.BulkUpload)
		protected.PUT("/employees/:id", middleware.RequireRoles("super_admin", "center_admin"), ctl.Employees.Update)
		protected.DELETE("/employees/:id", middleware.RequireRoles("super_admin", "center_admin"), ctl.Employees.Delete)

		protected.POST("/attendance/checkin", middleware.RequireRoles("super_admin", "center_admin", "operator"), ctl.Attendance.CheckIn)
		protected.POST("/attendance/checkout", middleware.RequireRoles("super_admin", "center_admin", "operator"), ctl.Attendance.CheckOut)
		protected.POST("/attendance/scan", middleware.RequireRoles("super_admin", "center_admin", "operator"), ctl.Attendance.Scan)
		protected.GET("/attendance", middleware.RequireRoles("super_admin", "center_admin", "operator"), ctl.Attendance.List)

		protected.GET("/reports/excel", middleware.RequireRoles("super_admin", "center_admin"), ctl.Reports.ExportExcel)
	}
}
