package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"
	"gorm.io/gorm"
)

// AssociationRepository is a repository for association operations
type OrganizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository creates a new Organization repository
func NewOrganizationRepository() *AssociationRepository {
	return &AssociationRepository{
		db: database.GetDB(),
	}
}

// FindByID finds a association by ID
func (r *OrganizationRepository) FindOrganizationByID(id uint) (*models.Organization, error) {
	var association models.Organization
	err := r.db.First(&association, id).Error
	if err != nil {
		return nil, err
	}
	return &association, nil
}