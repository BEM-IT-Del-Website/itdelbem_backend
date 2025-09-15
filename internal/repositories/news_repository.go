package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"
	"errors"

	"gorm.io/gorm"
)

// NewsRepository adalah repository untuk operasi terkait berita.
type NewsRepository struct {
	db *gorm.DB
}

// NewNewsRepository membuat instance news repository baru.
func NewNewsRepository() *NewsRepository {
	return &NewsRepository{
		db: database.GetDB(),
	}
}

// Create membuat item berita baru.
func (r *NewsRepository) Create(news *models.News) error {
	return r.db.Create(news).Error
}

// Update menyimpan perubahan pada item berita yang ada.
func (r *NewsRepository) Update(news *models.News) error {
	return r.db.Save(news).Error
}

// FindByID mencari item berita berdasarkan ID (hanya yang aktif).
func (r *NewsRepository) FindByID(id uint) (*models.News, error) {
	var news models.News
	err := r.db.First(&news, id).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}

// GetAllNews mengambil semua berita dengan pagination (hanya yang aktif).
func (r *NewsRepository) GetAllNews(limit, offset int) ([]models.News, int64, error) {
	var newsList []models.News
	var total int64

	// Query untuk menghitung total data yang aktif
	if err := r.db.Model(&models.News{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Query untuk mengambil data dengan limit, offset, dan pengurutan
	if err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&newsList).Error; err != nil {
		return nil, 0, err
	}

	return newsList, total, nil
}

// DeleteByID menghapus item berita berdasarkan ID (soft delete).
func (r *NewsRepository) DeleteByID(id uint) error {
	return r.db.Delete(&models.News{}, id).Error
}

// FindOnlyDeletedByID mencari item berita yang HANYA sudah di-soft-delete berdasarkan ID.
func (r *NewsRepository) FindOnlyDeletedByID(id uint) (*models.News, error) {
	var news models.News
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
func (r *NewsRepository) RestoreByID(id uint) (*models.News, error) {
	// 1. Cari record yang sudah dihapus
	deletedNews, err := r.FindOnlyDeletedByID(id)
	if err != nil {
		return nil, err
	}
	if deletedNews == nil {
		return nil, nil
	}

	// 2. Pulihkan record dengan mengupdate 'deleted_at' menjadi NULL
	if err := r.db.Unscoped().Model(&models.News{}).Where("id = ?", deletedNews.ID).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}

	// 3. Kembalikan record yang sudah dipulihkan (sekarang sudah aktif)
	return r.FindByID(deletedNews.ID)
}
