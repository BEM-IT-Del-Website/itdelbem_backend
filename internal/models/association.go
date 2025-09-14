package models

import (
	"time"

	"gorm.io/gorm"
)

// Association represents a student association (Himpunan).
type Association struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	ShortName string         `json:"short_name" gorm:"type:varchar(50);unique"`
	Vision    string         `json:"vision" gorm:"type:text"`
	Mission   string         `json:"mission" gorm:"type:text"`
	Values    string         `json:"values" gorm:"type:text"`
	Image     string         `json:"image" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Association) TableName() string {
	return "associations"
}

// AssociationManagement holds the board for an association in a specific period.
type AssociationManagement struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	AssociationID uint           `json:"association_id" gorm:"not null"`
	Association   *Association   `json:"association,omitempty" gorm:"foreignKey:AssociationID"`
	LeaderID      uint           `json:"leader_id" gorm:"not null"`
	Leader        *User          `json:"leader,omitempty" gorm:"foreignKey:LeaderID"`
	CoLeaderID    *uint          `json:"co_leader_id"`
	CoLeader      *User          `json:"co_leader,omitempty" gorm:"foreignKey:CoLeaderID"`
	Secretary1ID  *uint          `json:"secretary_1_id"`
	Secretary2ID  *uint          `json:"secretary_2_id"`
	Treasurer1ID  *uint          `json:"treasurer_1_id"`
	Treasurer2ID  *uint          `json:"treasurer_2_id"`
	Period        string         `json:"period" gorm:"type:varchar(20);not null"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (AssociationManagement) TableName() string {
	return "association_managements"
}

// AssociationWorkProgram represents a work program (proker) for an association.
type AssociationWorkProgram struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	AssociationID uint           `json:"association_id" gorm:"not null"`
	Association   *Association   `json:"association,omitempty" gorm:"foreignKey:AssociationID"`
	Name          string         `json:"name" gorm:"type:varchar(150);not null"`
	Description   string         `json:"description" gorm:"type:text"`
	Duration      string         `json:"duration" gorm:"type:varchar(100)"`
	Target        string         `json:"target" gorm:"type:varchar(255)"`
	Budget        float64        `json:"budget" gorm:"type:decimal(15,2)"`
	CoordinatorID *uint          `json:"coordinator_id"`
	Coordinator   *User          `json:"coordinator,omitempty" gorm:"foreignKey:CoordinatorID"`
	Status        string         `json:"status" gorm:"type:varchar(50)"`
	Objective     string         `json:"objective" gorm:"type:text"`
	PhotoURL      string         `json:"photo_url" gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (AssociationWorkProgram) TableName() string {
	return "association_work_programs"
}
