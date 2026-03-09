package models

import "time"

type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     *uint     `gorm:"index" json:"user_id"`
	Action     string    `gorm:"size:120;not null" json:"action"`
	EntityType string    `gorm:"size:60;not null" json:"entity_type"`
	EntityID   string    `gorm:"size:60;not null" json:"entity_id"`
	Details    string    `gorm:"type:text" json:"details"`
	CreatedAt  time.Time `json:"created_at"`
}
