package models

import (
	"time"

	"gorm.io/gorm"
)

type Center struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:120;uniqueIndex;not null" json:"name"`
	Location  string         `gorm:"size:255;not null" json:"location"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
