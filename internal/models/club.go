package models

import (
	"time"

	"gorm.io/gorm"
)

// Club represents a club in the system
type Club struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	ShortName string         `gorm:"not null" json:"shortname"`
	Vision    string         `gorm:"not null" json:"vision"`
	Mission   string         `gorm:"not null" json:"mission"`
	Value     string         `gorm:"not null" json:"value"`
	WorkPlan  string         `gorm:"not null" json:"workplan"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;uniqueIndex:idx_courses_code_deleted_at" json:"deleted_at,omitempty"`
}