package models

import (
	"time"

	"gorm.io/gorm"
)

// Organization represents a club in the system
type Organization struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CategoryID int            `form:"category_id" json:"category_id" gorm:"not null"`
	Category   *Category      `json:"category" gorm:"foreignKey:ID;references:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name       string         `form:"name" gorm:"not null" json:"name"`
	ShortName  string         `form:"short_name" gorm:"not null" json:"short_name"`
	Vision     string         `form:"vision" gorm:"not null" json:"vision"`
	Mission    string         `form:"mission" gorm:"not null" json:"mission"`
	Value      string         `form:"value" gorm:"not null" json:"value"`
	WorkPlan   string         `form:"workplan" gorm:"not null" json:"workplan"`
	Image      string         `form:"image" json:"image" gorm:"type:text"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index;uniqueIndex:idx_courses_code_deleted_at" json:"deleted_at,omitempty"`
}
