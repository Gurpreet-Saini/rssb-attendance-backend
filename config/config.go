package config

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"attendance-system/backend/models"
	"attendance-system/backend/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	AppEnv               string
	Port                 string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	DBSSLMode            string
	JWTSecret            string
	TokenDuration        time.Duration
	FrontendURL          string
	PasswordResetTTL     time.Duration
	SMTPHost             string
	SMTPPort             string
	SMTPUsername         string
	SMTPPassword         string
	SMTPFromEmail        string
	SMTPFromName         string
	DefaultAdminName     string
	DefaultAdminEmail    string
	DefaultAdminPassword string
}

func Load() Config {
	loadDotEnv(".env")

	return Config{
		AppEnv:               getEnv("APP_ENV", "development"),
		Port:                 getEnv("PORT", "8080"),
		DBHost:               getEnv("DB_HOST", "127.0.0.1"),
		DBPort:               getEnv("DB_PORT", "5433"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "postgres"),
		DBName:               getEnv("DB_NAME", "attendance_db"),
		DBSSLMode:            getEnv("DB_SSLMODE", "disable"),
		JWTSecret:            getEnv("JWT_SECRET", "change-me"),
		TokenDuration:        time.Duration(getEnvAsInt("JWT_DURATION_HOURS", 24)) * time.Hour,
		FrontendURL:          getEnv("FRONTEND_URL", "http://localhost:3001"),
		PasswordResetTTL:     time.Duration(getEnvAsInt("PASSWORD_RESET_TTL_MINUTES", 60)) * time.Minute,
		SMTPHost:             getEnv("SMTP_HOST", ""),
		SMTPPort:             getEnv("SMTP_PORT", "587"),
		SMTPUsername:         getEnv("SMTP_USERNAME", ""),
		SMTPPassword:         getEnv("SMTP_PASSWORD", ""),
		SMTPFromEmail:        getEnv("SMTP_FROM_EMAIL", ""),
		SMTPFromName:         getEnv("SMTP_FROM_NAME", "Attendance System"),
		DefaultAdminName:     getEnv("DEFAULT_ADMIN_NAME", "Super Admin"),
		DefaultAdminEmail:    getEnv("DEFAULT_ADMIN_EMAIL", "admin@example.com"),
		DefaultAdminPassword: getEnv("DEFAULT_ADMIN_PASSWORD", "Admin@123"),
	}
}

func ConnectDatabase(cfg Config) (*gorm.DB, error) {
	dsn := buildDSN(cfg)

	gormLogger := logger.Default.LogMode(logger.Info)
	if cfg.AppEnv == "production" {
		gormLogger = logger.Default.LogMode(logger.Warn)
	}

	var (
		db  *gorm.DB
		err error
	)

	for attempt := 1; attempt <= 10; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})
		if err == nil {
			var sqlDB *sql.DB
			sqlDB, err = db.DB()
			if err == nil {
				err = sqlDB.Ping()
			}
		}
		if err == nil {
			break
		}
		if attempt == 10 {
			return nil, err
		}
		log.Printf("database connection attempt %d failed: %v", attempt, err)
		time.Sleep(3 * time.Second)
	}

	if err := db.AutoMigrate(
		&models.Center{},
		&models.User{},
		&models.Employee{},
		&models.Attendance{},
		&models.AuditLog{},
	); err != nil {
		return nil, err
	}

	return db, nil
}

func EnsureDefaultAdmin(db *gorm.DB, cfg Config) error {
	var count int64
	if err := db.Model(&models.User{}).Where("role = ?", models.RoleSuperAdmin).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	password, err := utils.HashPassword(cfg.DefaultAdminPassword)
	if err != nil {
		return err
	}

	admin := models.User{
		Name:     cfg.DefaultAdminName,
		Email:    cfg.DefaultAdminEmail,
		Password: password,
		Role:     models.RoleSuperAdmin,
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Printf("seeded default super admin: %s", cfg.DefaultAdminEmail)
	return nil
}

func buildDSN(cfg Config) string {
	if databaseURL := getEnv("DATABASE_URL", ""); databaseURL != "" {
		if normalizedURL, ok := normalizeDatabaseURL(databaseURL, cfg.DBSSLMode); ok {
			return normalizedURL
		}
		log.Printf("invalid DATABASE_URL detected, falling back to DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME")
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Kolkata",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)
}

func normalizeDatabaseURL(databaseURL, fallbackSSLMode string) (string, bool) {
	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		return "", false
	}

	query := parsedURL.Query()
	if query.Get("sslmode") == "" {
		sslMode := fallbackSSLMode
		if sslMode == "" {
			sslMode = "require"
		}
		query.Set("sslmode", sslMode)
	}
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), true
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		_ = os.Setenv(key, value)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}
