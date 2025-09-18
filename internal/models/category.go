package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a club in the system
type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `form:"name" gorm:"not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;uniqueIndex:idx_courses_code_deleted_at" json:"deleted_at,omitempty"`
}
