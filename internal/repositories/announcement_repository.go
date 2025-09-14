package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"
	"gorm.io/gorm"
)

// AssociationRepository is a repository for association operations
type AnnouncementRepository struct {
	db *gorm.DB
}

// NewAssociationRepository creates a new association repository
func NewAnnouncementRepository() *AnnouncementRepository {
	return &AnnouncementRepository{
		db: database.GetDB(),
	}
}

// Create creates a new association
func (r *AnnouncementRepository) Create(announcement *models.Announcement) error {
	return r.db.Create(announcement).Error
}

// Update updates an existing association
func (r *AnnouncementRepository) Update(announcement *models.Announcement) error {
	return r.db.Save(announcement).Error
}

// FindByID finds a association by ID
func (r *AnnouncementRepository) FindByID(id uint) (*models.Announcement, error) {
	var announcement models.Announcement
	err := r.db.First(&announcement, id).Error
	if err != nil {
		return nil, err
	}
	return &announcement, nil
}

// FindByName finds a association by code
func (r *AnnouncementRepository) FindByName(code string) (*models.Announcement, error) {
	var announcement models.Announcement
	err := r.db.Where("code = ?", code).First(&announcement).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &announcement, nil
}

// FindAll finds all associations
func (r *AnnouncementRepository) GetAllAnnouncements(limit, offset int) ([]models.Announcement, int64, error) {
    var announcements []models.Announcement
    var total int64

    query := r.db.Model(&models.Announcement{})
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if err := query.Limit(limit).Offset(offset).Find(&announcements).Error; err != nil {
        return nil, 0, err
    }

    return announcements, total, nil
}


// DeleteByID deletes a association by ID
func (r *AnnouncementRepository) DeleteByID(id uint) error {
	// Use soft delete (don't use Unscoped())
	return r.db.Delete(&models.Announcement{}, id).Error
}

// FindDeletedByName finds a soft-deleted association by code
func (r *AnnouncementRepository) FindDeletedByName(code string) (*models.Announcement, error) {
	var announcement models.Announcement
	err := r.db.Unscoped().Where("code = ? AND deleted_at IS NOT NULL", code).First(&announcement).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &announcement, nil
}

// RestoreByName restores a soft-deleted association by code
func (r *AnnouncementRepository) RestoreByName(code string) (*models.Announcement, error) {
	// Find the deleted record
	deletedAnnouncement, err := r.FindDeletedByName(code)
	if err != nil {
		return nil, err
	}
	if deletedAnnouncement == nil {
		return nil, nil
	}
	
	// Restore the record
	if err := r.db.Unscoped().Model(&models.Announcement{}).Where("id = ?", deletedAnnouncement.ID).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}
	
	// Return the restored record
	return r.FindByID(deletedAnnouncement.ID)
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