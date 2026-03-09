package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	AttendancePresent = "present"
	AttendanceAbsent  = "absent"
	AttendancePartial = "partial"
)

type Attendance struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	EmployeeID uint           `gorm:"uniqueIndex:idx_attendance_employee_date;not null" json:"employee_id"`
	Employee   Employee       `json:"employee"`
	Date       time.Time      `gorm:"type:date;uniqueIndex:idx_attendance_employee_date;not null" json:"date"`
	TimeIn     *time.Time     `json:"time_in"`
	TimeOut    *time.Time     `json:"time_out"`
	Status     string         `gorm:"size:20;not null;index" json:"status"`
	CenterID   uint           `gorm:"index;not null" json:"center_id"`
	Center     Center         `json:"center"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
