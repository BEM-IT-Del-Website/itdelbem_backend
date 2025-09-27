package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"

	"gorm.io/gorm"
)

type RequestRepository struct {
	db *gorm.DB
}

func NewRequestRepository() *RequestRepository {
	return &RequestRepository{
		db: database.GetDB(),
	}
}

func (r *RequestRepository) Create(request *models.Request) error {
	return r.db.Create(request).Error
}

func (r *RequestRepository) Update(request *models.Request) error {
	return r.db.Save(request).Error
}

func (r *RequestRepository) FindByID(id uint) (*models.Request, error) {
	var request models.Request
	err := r.db.First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *RequestRepository) GetAllRequests(limit, offset int) ([]models.Request, int64, error) {
	var requests []models.Request
	var total int64

	query := r.db.Model(&models.Request{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

func (r *RequestRepository) DeleteByID(id uint) error {
	return r.db.Delete(&models.Request{}, id).Error
}
