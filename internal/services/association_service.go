package services

import (
	"gorm.io/gorm"
	"errors"

	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

// AssociationService is a service for association operations
type AssociationService struct {
	repository *repositories.AssociationRepository
	db *gorm.DB
}

// NewAssociationService creates a new association service
func NewAssociationService(db *gorm.DB) *AssociationService {
    return &AssociationService{
        repository: repositories.NewAssociationRepository(),
    }
}

// CreateAssociation creates a new association
func (s *AssociationService) CreateAssociation(association *models.Association) error {
	// Check if code exists (including soft-deleted)
	// exists, err := s.repository.CheckNameExists(association.Name, 0)
	// if err != nil {
	// 	return err
	// }

	// if exists {
	// 	// Try to find a soft-deleted association with this code
	// 	deletedAssociation, err := s.repository.FindDeletedByName(association.Name)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if deletedAssociation != nil {
	// 		// Restore the soft-deleted association with updated data
	// 		deletedAssociation.Name = association.Name
			
	// 		// Restore the association
	// 		restoredAssociation, err := s.repository.RestoreByName(association.Name)
	// 		if err != nil {
	// 			return err
	// 		}
			
	// 		// Update with new data
	// 		restoredAssociation.Name = association.Name
			
	// 		return s.repository.Update(restoredAssociation)
	// 	}
		
	// 	return errors.New("kode gedung sudah digunakan")
	// }

	// Create association
	return s.repository.Create(association)
}

// UpdateAssociation updates an existing association
func (s *AssociationService) UpdateAssociation(association *models.Association) error {
	// Check if association exists
	existingAssociation, err := s.repository.FindByID(association.ID)
	if err != nil {
		return err
	}
	if existingAssociation == nil {
		return errors.New("himpunan tidak ditemukan")
	}

	// Update association
	return s.repository.Update(association)
}

// GetAssociationByID gets a association by ID
func (s *AssociationService) GetAssociationByID(id uint) (*models.Association, error) {
	return s.repository.FindByID(id)
}

// GetAllAssociations gets all associations
func (s *AssociationService) GetAllAssociations(limit, offset int) ([]models.Association, int64, error) {
    return s.repository.GetAllAssociations(limit, offset)
}

// DeleteAssociation deletes a association
func (s *AssociationService) DeleteAssociation(id uint) error {
	// Check if association exists
	association, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if association == nil {
		return errors.New("gedung tidak ditemukan")
	}

	// Delete association (soft delete)
	return s.repository.DeleteByID(id)
}

// AssociationWithStats represents a association with additional statistics
type AssociationWithStats struct {
	Association  models.Association `json:"association"`
	RoomCount int64           `json:"room_count"`
}

// GetAssociationWithStats gets a association with its statistics
func (s *AssociationService) GetAssociationWithStats(id uint) (*AssociationWithStats, error) {
	// Get association
	association, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if association == nil {
		return nil, errors.New("gedung tidak ditemukan")
	}

	// Return association with stats
	return &AssociationWithStats{
		Association:  *association,
	}, nil
}

// GetAllAssociationsWithStats gets all associations with their statistics
// func (s *AssociationService) GetAllAssociationsWithStats() ([]AssociationWithStats, error) {
// 	// Get all associations
// 	associations, err := s.repository.Get()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Build response with stats
// 	result := make([]AssociationWithStats, len(associations))
// 	for i, association := range associations {
		
// 		result[i] = AssociationWithStats{
// 			Association:  association,
// 		}
// 	}

// 	return result, nil
// } 