package models

import (
	"time"

	"gorm.io/gorm"
)

type Request struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name" gorm:"size:255;not null"`
	Quantity         uint           `json:"quantity" gorm:"not null"`
	RequestPlan      string         `json:"request_plan" gorm:"not null"`
	ReturnPlan       string         `json:"return_plan" gorm:"not null"`
	RequesterID      int            `json:"requester_id" gorm:"not null"`
	ApproverID       *uint          `json:"approver_id" gorm:"default:null"`
	Requester        *User          `json:"requester,omitempty" gorm:"foreignKey:RequesterID"`
	Approver         *User          `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`
	Status           string         `json:"status" gorm:"type:enum('pending', 'approved', 'rejected');default:'pending';not null"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index;uniqueIndex:idx_courses_code_deleted_at" json:"deleted_at,omitempty"`
	OrganizationID   int            `json:"organization_id"`
	OrganizationName string         `json:"organization_name"`
}

func (Request) TableName() string {
	return "requests"
}
