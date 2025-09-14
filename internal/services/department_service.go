package services

import (
	"gorm.io/gorm"
	"errors"

	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

// DepartmentService is a service for department operations
type DepartmentService struct {
	repository *repositories.DepartmentRepository
	db *gorm.DB
}

// NewDepartmentService creates a new department service
func NewDepartmentService(db *gorm.DB) *DepartmentService {
    return &DepartmentService{
        repository: repositories.NewDepartmentRepository(),
    }
}

// CreateDepartment creates a new department
func (s *DepartmentService) CreateDepartment(department *models.Department) error {
	// Check if code exists (including soft-deleted)
	// exists, err := s.repository.CheckNameExists(department.Name, 0)
	// if err != nil {
	// 	return err
	// }

	// if exists {
	// 	// Try to find a soft-deleted department with this code
	// 	deletedDepartment, err := s.repository.FindDeletedByName(department.Name)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if deletedDepartment != nil {
	// 		// Restore the soft-deleted department with updated data
	// 		deletedDepartment.Name = department.Name
			
	// 		// Restore the department
	// 		restoredDepartment, err := s.repository.RestoreByName(department.Name)
	// 		if err != nil {
	// 			return err
	// 		}
			
	// 		// Update with new data
	// 		restoredDepartment.Name = department.Name
			
	// 		return s.repository.Update(restoredDepartment)
	// 	}
		
	// 	return errors.New("kode gedung sudah digunakan")
	// }

	// Create department
	return s.repository.Create(department)
}

// UpdateDepartment updates an existing department
func (s *DepartmentService) UpdateDepartment(department *models.Department) error {
	// Check if department exists
	existingDepartment, err := s.repository.FindByID(department.ID)
	if err != nil {
		return err
	}
	if existingDepartment == nil {
		return errors.New("himpunan tidak ditemukan")
	}

	// Update department
	return s.repository.Update(department)
}

// GetDepartmentByID gets a department by ID
func (s *DepartmentService) GetDepartmentByID(id uint) (*models.Department, error) {
	return s.repository.FindByID(id)
}

// GetAllDepartments gets all departments
func (s *DepartmentService) GetAllDepartments(limit, offset int) ([]models.Department, int64, error) {
    return s.repository.GetAllDepartments(limit, offset)
}

// DeleteDepartment deletes a department
func (s *DepartmentService) DeleteDepartment(id uint) error {
	// Check if department exists
	department, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if department == nil {
		return errors.New("gedung tidak ditemukan")
	}

	// Delete department (soft delete)
	return s.repository.DeleteByID(id)
}

// DepartmentWithStats represents a department with additional statistics
type DepartmentWithStats struct {
	Department  models.Department `json:"department"`
	RoomCount int64           `json:"room_count"`
}

// GetDepartmentWithStats gets a department with its statistics
func (s *DepartmentService) GetDepartmentWithStats(id uint) (*DepartmentWithStats, error) {
	// Get department
	department, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if department == nil {
		return nil, errors.New("gedung tidak ditemukan")
	}

	// Return department with stats
	return &DepartmentWithStats{
		Department:  *department,
	}, nil
}

// GetAllDepartmentsWithStats gets all departments with their statistics
// func (s *DepartmentService) GetAllDepartmentsWithStats() ([]DepartmentWithStats, error) {
// 	// Get all departments
// 	departments, err := s.repository.Get()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Build response with stats
// 	result := make([]DepartmentWithStats, len(departments))
// 	for i, department := range departments {
		
// 		result[i] = DepartmentWithStats{
// 			Department:  department,
// 		}
// 	}

// 	return result, nil
// } 