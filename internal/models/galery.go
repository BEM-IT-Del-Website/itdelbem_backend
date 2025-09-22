
package models

import (
	"time"

	"gorm.io/gorm"
)

type Galery struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:255;not null"`
	Content   string         `json:"content" gorm:"not null"`
	ImageURL  string         `json:"image_url" gorm:"type:varchar(255)"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;uniqueIndex:idx_courses_code_deleted_at" json:"deleted_at,omitempty"`
}

func (Galery) TableName() string {
	return "galery"
}