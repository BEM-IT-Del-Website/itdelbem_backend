package repositories

import (
	"log"
	"errors"

	"bem_be/internal/database"
	"bem_be/internal/models"

	"gorm.io/gorm"
)

// StudentRepository handles database operations for students
type StudentRepository struct {
	db *gorm.DB
}

// NewStudentRepository creates a new student repository
func NewStudentRepository() *StudentRepository {
	return &StudentRepository{
		db: database.GetDB(),
	}
}

func (r *StudentRepository) Update(student *models.Student) error {
	if student.ID == 0 {
		return errors.New("student ID is required")
	}
	return r.db.Save(student).Error
}

// FindAll returns all students from the database with optional search and filters
func (r *StudentRepository) FindAll(limit, offset int, search, studyProgram string, yearEnrolled int) ([]models.Student, int64, error) {
    var students []models.Student
    var total int64

    query := r.db.Model(&models.Student{})

    // filter by search (di full_name, study_program, year_enrolled)
    if search != "" {
        likeSearch := "%" + search + "%"
        query = query.Where(
            r.db.Where("full_name LIKE ?", likeSearch).
                Or("study_program LIKE ?", likeSearch).
                Or("CAST(year_enrolled AS CHAR) LIKE ?", likeSearch),
        )
    }

    // filter by study program (pakai LIKE biar fleksibel)
    if studyProgram != "" {
        query = query.Where("study_program LIKE ?", "%"+studyProgram+"%")
    }

    // filter by year enrolled
    if yearEnrolled > 0 {
        query = query.Where("year_enrolled = ?", yearEnrolled)
    }

    // hitung total sesuai filter
    query.Count(&total)

    // ambil data
    result := query.
        Order("year_enrolled ASC").
        Order("study_program ASC").
        Order("nim ASC").
        Limit(limit).
        Offset(offset).
        Find(&students)

    return students, total, result.Error
}
// FindByID returns a student by ID
func (r *StudentRepository) FindByID(id uint) (*models.Student, error) {
	var student models.Student
	result := r.db.First(&student, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &student, nil
}

// FindByNIM returns a student by NIM
func (r *StudentRepository) FindByNIM(nim string) (*models.Student, error) {
	var student models.Student
	result := r.db.Where("nim = ?", nim).First(&student)
	if result.Error != nil {
		return nil, result.Error
	}
	return &student, nil
}

// FindByUserID returns a student by external UserID from campus
func (r *StudentRepository) FindByUserID(userID int) (*models.Student, error) {
	var student models.Student
	result := r.db.Where("user_id = ?", userID).First(&student)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &student, nil
}

// UpsertMany creates or updates multiple students
func (r *StudentRepository) UpsertMany(students []models.Student) error {
	if len(students) == 0 {
		return nil
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, student := range students {
		// Try to find existing student by DimID (from external system)
		var existingStudent models.Student
		result := tx.Where("dim_id = ?", student.DimID).First(&existingStudent)

		if result.Error == nil {
			// Check if the student ID is going to change
			oldID := existingStudent.ID

			// Update existing student
			student.ID = existingStudent.ID
			student.CreatedAt = existingStudent.CreatedAt

			if err := tx.Save(&student).Error; err != nil {
				tx.Rollback()
				return err
			}

			// Update student_to_groups rows if the student ID changed but UserID remains the same
			// This maintains group membership connections when student IDs change
			if oldID != student.ID && existingStudent.UserID == student.UserID {
				if err := tx.Exec(
					"UPDATE student_to_groups SET student_id = ? WHERE student_id = ? AND user_id = ?",
					student.ID, oldID, student.UserID,
				).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		} else {
			// Create new student
			if err := tx.Create(&student).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	log.Printf("Upserted %d students", len(students))
	return tx.Commit().Error
}

func (r *StudentRepository) FindByExternalUserUsername(username string) (*models.Student, error) {
	var student models.Student
	if err := r.db.Where("user_name = ?", username).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

