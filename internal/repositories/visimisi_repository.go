package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"
	"errors"

	"gorm.io/gorm"
)

// VisiMisiRepository adalah repository untuk operasi terkait berita.
type VisiMisiRepository struct {
	db *gorm.DB
}

// NewVisiMisiRepository membuat instance news repository baru.
func NewVisiMisiRepository() *VisiMisiRepository {
	return &VisiMisiRepository{
		db: database.GetDB(),
	}
}

// Create membuat item berita baru.
func (r *VisiMisiRepository) Create(news *models.Period) error {
	return r.db.Create(news).Error
}

// Update menyimpan perubahan pada item berita yang ada.
func (r *VisiMisiRepository) Update(news *models.Period) error {
	return r.db.Save(news).Error
}

// FindByID mencari item berita berdasarkan ID (hanya yang aktif).
func (r *VisiMisiRepository) FindByID(id uint) (*models.Period, error) {
	var news models.Period
	err := r.db.First(&news, id).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}

// GetAllVisiMisi mengambil semua berita dengan pagination (hanya yang aktif).
func (r *VisiMisiRepository) GetAllVisiMisi(limit, offset int) ([]models.Period, int64, error) {
	var newsList []models.Period
	var total int64

	// Query untuk menghitung total data yang aktif
	if err := r.db.Model(&models.Period{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Query untuk mengambil data dengan limit, offset, dan pengurutan
	if err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&newsList).Error; err != nil {
		return nil, 0, err
	}

	return newsList, total, nil
}

// DeleteByID menghapus item berita berdasarkan ID (soft delete).
func (r *VisiMisiRepository) DeleteByID(id uint) error {
	return r.db.Delete(&models.Period{}, id).Error
}

// FindOnlyDeletedByID mencari item berita yang HANYA sudah di-soft-delete berdasarkan ID.
func (r *VisiMisiRepository) FindOnlyDeletedByID(id uint) (*models.Period, error) {
	var news models.Period
	err := r.db.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&news).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &news, nil
}

// RestoreByID memulihkan item berita yang di-soft-delete berdasarkan ID.
func (r *VisiMisiRepository) RestoreByID(id uint) (*models.Period, error) {
	// 1. Cari record yang sudah dihapus
	deletedVisiMisi, err := r.FindOnlyDeletedByID(id)
	if err != nil {
		return nil, err
	}
	if deletedVisiMisi == nil {
		return nil, nil
	}

	// 2. Pulihkan record dengan mengupdate 'deleted_at' menjadi NULL
	if err := r.db.Unscoped().Model(&models.Period{}).Where("id = ?", deletedVisiMisi.ID).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}

	// 3. Kembalikan record yang sudah dipulihkan (sekarang sudah aktif)
	return r.FindByID(deletedVisiMisi.ID)
}
