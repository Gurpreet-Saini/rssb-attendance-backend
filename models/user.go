package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	RoleSuperAdmin  = "super_admin"
	RoleCenterAdmin = "center_admin"
	RoleOperator    = "operator"
)

type User struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	Name                 string         `gorm:"size:120;not null" json:"name"`
	Email                string         `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Password             string         `gorm:"not null" json:"-"`
	Role                 string         `gorm:"size:30;not null;index" json:"role"`
	CenterID             *uint          `gorm:"index" json:"center_id"`
	Center               *Center        `json:"center,omitempty"`
	PasswordResetToken   *string        `gorm:"size:128;index" json:"-"`
	PasswordResetExpires *time.Time     `json:"-"`
	CreatedAt            time.Time      `json:"created_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
