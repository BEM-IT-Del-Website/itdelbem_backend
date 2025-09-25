package models

import (
	"time"

	"gorm.io/gorm"
)

// BEM represents the main student executive board for a specific period.
type BEM struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	LeaderID     uint           `json:"leader_id"`
	Leader       *Student       `json:"leader" gorm:"foreignKey:ID;references:LeaderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CoLeaderID   uint           `json:"coleader_id"`
	CoLeader     *Student       `json:"coleader" gorm:"foreignKey:ID;references:CoLeaderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Secretary1ID uint           `json:"secretary1_id"`
	Secretary1   *Student       `json:"secretary1" gorm:"foreignKey:ID;references:Secretary1ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Secretary2ID uint           `json:"secretary2_id"`
	Secretary2   *Student       `json:"secretary2" gorm:"foreignKey:ID;references:Secretary2ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Treasurer1ID uint           `json:"treasurer1_id"`
	Treasurer1   *Student       `json:"treasurer1" gorm:"foreignKey:ID;references:Treasurer1ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Treasurer2ID uint           `json:"treasurer2_id"`
	Treasurer2   *Student       `json:"treasurer2" gorm:"foreignKey:ID;references:Treasurer2ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Period       string         `json:"period" gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (BEM) TableName() string {
	return "bems"
}
