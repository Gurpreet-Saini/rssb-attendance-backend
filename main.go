package main

import (
	"log"
	"time"

	"attendance-system/backend/config"
	"attendance-system/backend/controllers"
	"attendance-system/backend/middleware"
	"attendance-system/backend/repositories"
	"attendance-system/backend/routes"
	"attendance-system/backend/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	if err := config.EnsureDefaultAdmin(db, cfg); err != nil {
		log.Fatalf("default admin setup failed: %v", err)
	}

	centerRepo := repositories.NewCenterRepository(db)
	userRepo := repositories.NewUserRepository(db)
	employeeRepo := repositories.NewEmployeeRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)
	auditRepo := repositories.NewAuditRepository(db)
	emailService := services.NewEmailService(cfg)

	authService := services.NewAuthService(userRepo, emailService, cfg)
	centerService := services.NewCenterService(centerRepo)
	userService := services.NewUserService(userRepo, auditRepo, emailService, cfg)
	employeeService := services.NewEmployeeService(employeeRepo, auditRepo)
	attendanceService := services.NewAttendanceService(attendanceRepo, employeeRepo)
	reportService := services.NewReportService(attendanceService)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:  cfg.CORSAllowedOrigins,
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Content-Length", "Authorization", "Accept"},
		ExposeHeaders: []string{"Content-Disposition", "File-Name"},
		MaxAge:        12 * time.Hour,
	}))
	router.Use(middleware.RateLimitMiddleware(120))

	routes.Register(router, cfg.JWTSecret, routes.Controllers{
		Auth:       controllers.NewAuthController(authService),
		Centers:    controllers.NewCenterController(centerService),
		Users:      controllers.NewUserController(userService),
		Employees:  controllers.NewEmployeeController(employeeService),
		Attendance: controllers.NewAttendanceController(attendanceService),
		Reports:    controllers.NewReportController(reportService),
	})

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
