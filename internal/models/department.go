package models

import (
	"time"

	"gorm.io/gorm"
)

// Department represents a department entity.
type Department struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	ShortName string         `json:"short_name" gorm:"type:varchar(50);unique"`
	Vision    string         `json:"vision" gorm:"type:text"`
	Mission   string         `json:"mission" gorm:"type:text"`
	Values    string         `json:"values" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Department) TableName() string {
	return "departments"
}

// DepartmentManagement holds the board for a department in a specific period.
type DepartmentManagement struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	DepartmentID uint           `json:"department_id" gorm:"not null"`
	Department   *Department    `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	LeaderID     uint           `json:"leader_id" gorm:"not null"`
	Leader       *User          `json:"leader,omitempty" gorm:"foreignKey:LeaderID"`
	CoLeaderID   *uint          `json:"co_leader_id"`
	CoLeader     *User          `json:"co_leader,omitempty" gorm:"foreignKey:CoLeaderID"`
	Secretary1ID *uint          `json:"secretary_1_id"`
	Secretary2ID *uint          `json:"secretary_2_id"`
	Treasurer1ID *uint          `json:"treasurer_1_id"`
	Treasurer2ID *uint          `json:"treasurer_2_id"`
	Period       string         `json:"period" gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (DepartmentManagement) TableName() string {
	return "department_managements"
}

// DepartmentWorkProgram represents a work program (proker) for a department.
type DepartmentWorkProgram struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	DepartmentID  uint           `json:"department_id" gorm:"not null"`
	Department    *Department    `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
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

func (DepartmentWorkProgram) TableName() string {
	return "department_work_programs"
}
