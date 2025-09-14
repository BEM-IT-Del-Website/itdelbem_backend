package repositories

import (
	"bem_be/internal/models"

	"gorm.io/gorm"
)

type UkmRepository struct {
	db *gorm.DB
}

func NewUkmRepository(db *gorm.DB) *UkmRepository {
	return &UkmRepository{db: db}
}

// CreateUkm akan menyimpan data UKM baru ke database
func (r *UkmRepository) CreateUkm(ukm *models.Ukm) error {
	return r.db.Create(ukm).Error
}

// Tambahkan fungsi lain di sini nanti (GetAll, GetByID, etc.)
