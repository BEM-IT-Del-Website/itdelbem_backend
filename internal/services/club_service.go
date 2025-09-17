package services

import (
	"gorm.io/gorm"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

// ClubService is a service for club operations
type ClubService struct {
	repository *repositories.ClubRepository
	db *gorm.DB
}

// NewClubService creates a new club service
func NewClubService(db *gorm.DB) *ClubService {
    return &ClubService{
        repository: repositories.NewClubRepository(),
    }
}

// CreateClub creates a new club
func (s *ClubService) CreateClub(association *models.Club, file *multipart.FileHeader) error {
	// bikin folder kalau belum ada
	if err := os.MkdirAll("uploads/departments", os.ModePerm); err != nil {
		return err
	}

	// nama file unik
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filepath := "uploads/departments/" + filename

	// simpan file
	if err := saveUploadedFile(file, filepath); err != nil {
		return err
	}

	// simpan path/filename ke struct
	association.Image = filename

	// simpan ke DB
	return s.repository.Create(association)
}

// UpdateClub updates an existing club
func (s *ClubService) UpdateClub(club *models.Club) error {
	// Check if club exists
	existingClub, err := s.repository.FindByID(club.ID)
	if err != nil {
		return err
	}
	if existingClub == nil {
		return errors.New("himpunan tidak ditemukan")
	}

	// Update club
	return s.repository.Update(club)
}

// GetClubByID gets a club by ID
func (s *ClubService) GetClubByID(id uint) (*models.Club, error) {
	return s.repository.FindByID(id)
}

// GetAllClubs gets all clubs
func (s *ClubService) GetAllClubs(limit, offset int) ([]models.Club, int64, error) {
    return s.repository.GetAllClubs(limit, offset)
}

func (s *ClubService) GetAllClubsGuest() ([]models.Club, error) {
    return s.repository.GetAllClubsGuest()
}

// DeleteClub deletes a club
func (s *ClubService) DeleteClub(id uint) error {
	// Check if club exists
	club, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if club == nil {
		return errors.New("gedung tidak ditemukan")
	}

	// Delete club (soft delete)
	return s.repository.DeleteByID(id)
}

// ClubWithStats represents a club with additional statistics
type ClubWithStats struct {
	Club  models.Club `json:"club"`
	RoomCount int64           `json:"room_count"`
}

// GetClubWithStats gets a club with its statistics
func (s *ClubService) GetClubWithStats(id uint) (*ClubWithStats, error) {
	// Get club
	club, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if club == nil {
		return nil, errors.New("gedung tidak ditemukan")
	}

	// Return club with stats
	return &ClubWithStats{
		Club:  *club,
	}, nil
}

// GetAllClubsWithStats gets all clubs with their statistics
// func (s *ClubService) GetAllClubsWithStats() ([]ClubWithStats, error) {
// 	// Get all clubs
// 	clubs, err := s.repository.Get()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Build response with stats
// 	result := make([]ClubWithStats, len(clubs))
// 	for i, club := range clubs {
		
// 		result[i] = ClubWithStats{
// 			Club:  club,
// 		}
// 	}

// 	return result, nil
// } 