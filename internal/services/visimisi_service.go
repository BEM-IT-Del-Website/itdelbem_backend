package services

import (
	"bem_be/internal/models"
	"bem_be/internal/repositories"
	"errors"

	"gorm.io/gorm"
)

// VisiMisiService adalah service untuk operasi berita.
type VisiMisiService struct {
	repository *repositories.VisiMisiRepository
}

// NewVisiMisiService membuat service berita baru.
func NewVisiMisiService(db *gorm.DB) *VisiMisiService {
	return &VisiMisiService{
		repository: repositories.NewVisiMisiRepository(),
	}
}

// CreateVisiMisi membuat berita baru.
func (s *VisiMisiService) CreateVisiMisi(visimisi *models.Period) error {
	if visimisi.Vision == "" || visimisi.Mission == "" {
		return errors.New("judul dan konten tidak boleh kosong")
	}
	return s.repository.Create(visimisi)
}

// UpdateVisiMisi memperbarui berita yang ada.
func (s *VisiMisiService) UpdateVisiMisi(visimisi *models.Period) error {
	return s.repository.Update(visimisi)
}

// GetVisiMisiByID mendapatkan berita berdasarkan ID.
func (s *VisiMisiService) GetVisiMisiByID(id uint) (*models.Period, error) {
	visimisi, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("berita tidak ditemukan")
		}
		return nil, err
	}
	return visimisi, nil
}

// GetAllVisiMisi mendapatkan semua berita dengan pagination.
func (s *VisiMisiService) GetAllVisiMisi(limit, offset int) ([]models.Period, int64, error) {
	return s.repository.GetAllVisiMisi(limit, offset)
}

// DeleteVisiMisi menghapus sebuah berita.
func (s *VisiMisiService) DeleteVisiMisi(id uint) error {
	_, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("berita yang akan dihapus tidak ditemukan")
		}
		return err
	}
	return s.repository.DeleteByID(id)
}

// RestoreVisiMisi memulihkan berita dan mengembalikan data yang telah dipulihkan.
func (s *VisiMisiService) RestoreVisiMisi(id uint) (*models.Period, error) {
	restoredVisiMisi, err := s.repository.RestoreByID(id)
	if err != nil {
		return nil, err
	}
	if restoredVisiMisi == nil {
		return nil, errors.New("berita tidak ditemukan atau sudah aktif (tidak dihapus)")
	}
	return restoredVisiMisi, nil
}
