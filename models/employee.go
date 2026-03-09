package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:120;not null;index" json:"name"`
	EmployeeID  string         `gorm:"size:50;uniqueIndex;not null" json:"employee_id"`
	BadgeNumber string         `gorm:"size:50;uniqueIndex;not null" json:"badge_number"`
	Barcode     string         `gorm:"size:80;uniqueIndex;not null" json:"barcode"`
	CenterID    uint           `gorm:"index;not null" json:"center_id"`
	Center      Center         `json:"center"`
	Designation string         `gorm:"size:120;not null" json:"designation"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
