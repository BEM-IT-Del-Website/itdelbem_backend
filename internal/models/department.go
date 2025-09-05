package models

import (
	"time"

	"gorm.io/gorm"
)

// Course represents a course in the system
type Course struct {
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

// Department represents a department/study program
type Department struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex:idx_departments_name_deleted_at;not null" json:"name"`
	FacultyID uint           `gorm:"not null" json:"faculty_id"`
	Faculty   Faculty        `gorm:"foreignKey:FacultyID" json:"faculty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;uniqueIndex:idx_departments_name_deleted_at" json:"deleted_at,omitempty"`
}

// Faculty struct is defined in faculty.go - removed duplicate declaration
