package services

import (
	"errors"

	"gorm.io/gorm"

	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

type GaleryService struct {
	repository *repositories.GaleryRepository
	db         *gorm.DB
}

func NewGaleryService(db *gorm.DB) *GaleryService {
	return &GaleryService{
		repository: repositories.NewGaleryRepository(),
	}
}

func (s *GaleryService) CreateGalery(galery *models.Galery) error {
	return s.repository.Create(galery)
}

func (s *GaleryService) UpdateGalery(galery *models.Galery) error {
	existingGalery, err := s.repository.FindByID(galery.ID)
	if err != nil {
		return err
	}
	if existingGalery == nil {
		return errors.New("gambar tidak ditemukan")
	}
	return s.repository.Update(galery)
}

func (s *GaleryService) GetGaleryByID(id uint) (*models.Galery, error) {
	return s.repository.FindByID(id)
}

func (s *GaleryService) GetAllGalerys(limit, offset int) ([]models.Galery, int64, error) {
	return s.repository.GetAllGalerys(limit, offset)
}

func (s *GaleryService) DeleteGalery(id uint) error {
	galery, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if galery == nil {
		return errors.New("gambar tidak ditemukan")
	}
	return s.repository.DeleteByID(id)
}

type GaleryWithStats struct {
	Galery    models.Galery `json:"galery"`
	RoomCount int64         `json:"room_count"`
}

func (s *GaleryService) GetGaleryWithStats(id uint) (*GaleryWithStats, error) {
	galery, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if galery == nil {
		return nil, errors.New("gedung tidak ditemukan")
	}
	return &GaleryWithStats{
		Galery: *galery,
	}, nil
}
