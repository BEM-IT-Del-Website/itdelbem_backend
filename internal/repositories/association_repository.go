package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"
	"gorm.io/gorm"
)

// AssociationRepository is a repository for association operations
type AssociationRepository struct {
	db *gorm.DB
}

// NewAssociationRepository creates a new association repository
func NewAssociationRepository() *AssociationRepository {
	return &AssociationRepository{
		db: database.GetDB(),
	}
}

// Create creates a new association
func (r *AssociationRepository) Create(association *models.Association) error {
	return r.db.Create(association).Error
}

// Update updates an existing association
func (r *AssociationRepository) Update(association *models.Association) error {
	return r.db.Save(association).Error
}

// FindByID finds a association by ID
func (r *AssociationRepository) FindByID(id uint) (*models.Association, error) {
	var association models.Association
	err := r.db.First(&association, id).Error
	if err != nil {
		return nil, err
	}
	return &association, nil
}

// FindByName finds a association by code
func (r *AssociationRepository) FindByName(code string) (*models.Association, error) {
	var association models.Association
	err := r.db.Where("code = ?", code).First(&association).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &association, nil
}

// GetAllAssociations returns all associations from the database with optional search filter
func (r *AssociationRepository) GetAllAssociations(limit, offset int, search string) ([]models.Association, int64, error) {
    var associations []models.Association
    var total int64

    query := r.db.Model(&models.Association{})

    if search != "" {
        likeSearch := "%" + search + "%"
        query = query.Where("name LIKE ?", likeSearch)
    }

    query.Count(&total)

    result := query.
        Order("name ASC").
        Limit(limit).
        Offset(offset).
        Find(&associations)

    return associations, total, result.Error
}


// DeleteByID deletes a association by ID
func (r *AssociationRepository) DeleteByID(id uint) error {
	// Use soft delete (don't use Unscoped())
	return r.db.Delete(&models.Association{}, id).Error
}

// FindDeletedByName finds a soft-deleted association by code
func (r *AssociationRepository) FindDeletedByName(code string) (*models.Association, error) {
	var association models.Association
	err := r.db.Unscoped().Where("code = ? AND deleted_at IS NOT NULL", code).First(&association).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &association, nil
}

// RestoreByName restores a soft-deleted association by code
func (r *AssociationRepository) RestoreByName(code string) (*models.Association, error) {
	// Find the deleted record
	deletedAssociation, err := r.FindDeletedByName(code)
	if err != nil {
		return nil, err
	}
	if deletedAssociation == nil {
		return nil, nil
	}
	
	// Restore the record
	if err := r.db.Unscoped().Model(&models.Association{}).Where("id = ?", deletedAssociation.ID).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}
	
	// Return the restored record
	return r.FindByID(deletedAssociation.ID)
}

// // CheckNameExists checks if a code exists, including soft-deleted records
// func (r *AssociationRepository) CheckNameExists(code string, excludeID uint) (bool, error) {
// 	var count int64
// 	query := r.db.Unscoped().Model(&models.Association{}).Where("code = ?", code)
	
// 	// Exclude the current record if updating
// 	if excludeID > 0 {
// 		query = query.Where("id != ?", excludeID)
// 	}
	
// 	err := query.Count(&count).Error
// 	if err != nil {
// 		return false, err
// 	}
	
// 	return count > 0, nil
// } 

func (r *AssociationRepository) GetAllAssociationsGuest() ([]models.Association, error) {
    var associations []models.Association
    err := r.db.Find(&associations).Error
    return associations, err
}
