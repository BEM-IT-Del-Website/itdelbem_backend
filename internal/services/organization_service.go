package services

import (
	"errors"

	"gorm.io/gorm"

	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

// OrganizationService is a service for organization operations
type OrganizationService struct {
	repository *repositories.OrganizationRepository
	db *gorm.DB
}

func NewOrganizationService(db *gorm.DB) *OrganizationService {
    return &OrganizationService{
        // repository: repositories.NewOrganizationRepository(),
    }
}

type OrganizationWithStats struct {
	Organization models.Organization `json:"organization"`
}

// GetOrganizationByID gets a organization by ID
func (s *OrganizationService) GetOrganizationByID(id uint) (*models.Organization, error) {
	return s.repository.FindOrganizationByID(id)
}

// GetOrganizationWithStats gets a organization with its statistics
func (s *OrganizationService) GetOrganizationWithStats(id uint) (*OrganizationWithStats, error) {
	// Get organization
	organization, err := s.repository.FindOrganizationByID(id)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, errors.New("organisasi tidak ditemukan")
	}

	// Return organization with stats
	return &OrganizationWithStats{
		Organization: *organization,
	}, nil
}
