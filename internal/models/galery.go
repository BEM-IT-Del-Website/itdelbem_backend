package models

import (
	"time"

	"gorm.io/gorm"
)

type Galery struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:255;not null"`
	Content   string         `json:"content" gorm:"not null"`
	Photo     string         `json:"photo" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;uniqueIndex:idx_courses_code_deleted_at" json:"deleted_at,omitempty"`
}
