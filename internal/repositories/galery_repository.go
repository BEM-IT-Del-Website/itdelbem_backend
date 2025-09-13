package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"

	"gorm.io/gorm"
)

type GaleryRepository struct {
	db *gorm.DB
}

func NewGaleryRepository() *GaleryRepository {
	return &GaleryRepository{
		db: database.GetDB(),
	}
}

func (r *GaleryRepository) Create(galery *models.Galery) error {
	return r.db.Create(galery).Error
}

func (r *GaleryRepository) Update(galery *models.Galery) error {
	return r.db.Save(galery).Error
}

func (r *GaleryRepository) FindByID(id uint) (*models.Galery, error) {
	var galery models.Galery
	err := r.db.First(&galery, id).Error
	if err != nil {
		return nil, err
	}
	return &galery, nil
}

func (r *GaleryRepository) GetAllGalerys(limit, offset int) ([]models.Galery, int64, error) {
	var galerys []models.Galery
	var total int64

	query := r.db.Model(&models.Galery{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&galerys).Error; err != nil {
		return nil, 0, err
	}

	return galerys, total, nil
}

func (r *GaleryRepository) DeleteByID(id uint) error {
	return r.db.Delete(&models.Galery{}, id).Error
}
