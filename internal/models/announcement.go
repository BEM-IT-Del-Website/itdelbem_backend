package models

import (
    "time"
    "gorm.io/gorm"
)

type Announcement struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Title     string         `json:"title" gorm:"size:255;not null"`
    Content   string         `json:"content" gorm:"type:text;not null"`
    FilePath  string         `json:"file_path,omitempty" gorm:"type:varchar(255)"`
    AuthorID  uint           `json:"author_id" gorm:"not null"`
    Author    *User          `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
    StartDate *time.Time     `json:"start_date,omitempty"`
    EndDate   *time.Time     `json:"end_date,omitempty"`
    CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
