package services

import (
	"errors"

	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

// CommunityService is a service for building operations
type CommunityService struct {
	repository *repositories.CommunityRepository
}

// NewCommunityService creates a new building service
func NewCommunityService() *CommunityService {
	return &CommunityService{
		repository: repositories.NewCommunityRepository(),
	}
}

// CreateCommunity creates a new building
func (s *CommunityService) CreateCommunity(building *models.Community) error {
	// Check if code exists (including soft-deleted)
	exists, err := s.repository.CheckCodeExists(building.Code, 0)
	if err != nil {
		return err
	}

	if exists {
		// Try to find a soft-deleted building with this code
		deletedCommunity, err := s.repository.FindDeletedByCode(building.Code)
		if err != nil {
			return err
		}

		if deletedCommunity != nil {
			// Restore the soft-deleted building with updated data
			deletedCommunity.Name = building.Name
			deletedCommunity.Floors = building.Floors
			deletedCommunity.Description = building.Description
			
			// Restore the building
			restoredCommunity, err := s.repository.RestoreByCode(building.Code)
			if err != nil {
				return err
			}
			
			// Update with new data
			restoredCommunity.Name = building.Name
			restoredCommunity.Floors = building.Floors
			restoredCommunity.Description = building.Description
			
			return s.repository.Update(restoredCommunity)
		}
		
		return errors.New("kode gedung sudah digunakan")
	}

	// Create building
	return s.repository.Create(building)
}

// UpdateCommunity updates an existing building
func (s *CommunityService) UpdateCommunity(building *models.Community) error {
	// Check if building exists
	existingCommunity, err := s.repository.FindByID(building.ID)
	if err != nil {
		return err
	}
	if existingCommunity == nil {
		return errors.New("gedung tidak ditemukan")
	}

	// If code is changed, check if new code already exists
	if building.Code != existingCommunity.Code {
		exists, err := s.repository.CheckCodeExists(building.Code, building.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("kode gedung sudah digunakan")
		}
	}

	// Update building
	return s.repository.Update(building)
}

// GetCommunityByID gets a building by ID
func (s *CommunityService) GetCommunityByID(id uint) (*models.Community, error) {
	return s.repository.FindByID(id)
}

// GetAllCommunitys gets all buildings
func (s *CommunityService) GetAllCommunitys() ([]models.Community, error) {
	return s.repository.FindAll()
}

// DeleteCommunity deletes a building
func (s *CommunityService) DeleteCommunity(id uint) error {
	// Check if building exists
	building, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if building == nil {
		return errors.New("gedung tidak ditemukan")
	}

	// Check if there are any associated rooms
	roomCount, err := s.repository.CountRooms(id)
	if err != nil {
		return err
	}
	if roomCount > 0 {
		return errors.New("tidak dapat menghapus gedung yang memiliki ruangan")
	}

	// Delete building (soft delete)
	return s.repository.DeleteByID(id)
}

// CommunityWithStats represents a building with additional statistics
type CommunityWithStats struct {
	Community  models.Community `json:"building"`
	RoomCount int64           `json:"room_count"`
}

// GetCommunityWithStats gets a building with its statistics
func (s *CommunityService) GetCommunityWithStats(id uint) (*CommunityWithStats, error) {
	// Get building
	building, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if building == nil {
		return nil, errors.New("gedung tidak ditemukan")
	}

	// Count rooms
	roomCount, err := s.repository.CountRooms(id)
	if err != nil {
		return nil, err
	}

	// Return building with stats
	return &CommunityWithStats{
		Community:  *building,
		RoomCount: roomCount,
	}, nil
}

// GetAllCommunitysWithStats gets all buildings with their statistics
func (s *CommunityService) GetAllCommunitysWithStats() ([]CommunityWithStats, error) {
	// Get all buildings
	buildings, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	// Build response with stats
	result := make([]CommunityWithStats, len(buildings))
	for i, building := range buildings {
		// Count rooms for each building
		roomCount, err := s.repository.CountRooms(building.ID)
		if err != nil {
			return nil, err
		}

		result[i] = CommunityWithStats{
			Community:  building,
			RoomCount: roomCount,
		}
	}

	return result, nil
} 