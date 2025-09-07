package models

import (
	"time"

	"gorm.io/gorm"
)

// BEM represents the main student executive board for a specific period.
type BEM struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	LeaderID     uint           `json:"leader_id" gorm:"not null"`
	Leader       *User          `json:"leader,omitempty" gorm:"foreignKey:LeaderID"`
	CoLeaderID   *uint          `json:"co_leader_id"`
	CoLeader     *User          `json:"co_leader,omitempty" gorm:"foreignKey:CoLeaderID"`
	Secretary1ID *uint          `json:"secretary_1_id"`
	Secretary2ID *uint          `json:"secretary_2_id"`
	Treasurer1ID *uint          `json:"treasurer_1_id"`
	Treasurer2ID *uint          `json:"treasurer_2_id"`
	Period       string         `json:"period" gorm:"type:varchar(20);not null;unique"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (BEM) TableName() string {
	return "bems"
}
