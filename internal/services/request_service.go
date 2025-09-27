package services

import (
	"bem_be/internal/models"
	"bem_be/internal/repositories"
	"errors"

	"gorm.io/gorm"
)

type RequestService struct {
	repository  *repositories.RequestRepository
	studentRepo *repositories.StudentRepository
	db          *gorm.DB
}

func NewRequestService(db *gorm.DB) *RequestService {
	return &RequestService{
		repository:  repositories.NewRequestRepository(),
		studentRepo: repositories.NewStudentRepository(),
	}
}

func (s *RequestService) GetStudentByUserID(userID int) (*models.Student, error) {
	return s.studentRepo.FindByUserID(userID)
}

func (s *RequestService) CreateRequest(request *models.Request) error {
	return s.repository.Create(request)
}

func (s *RequestService) UpdateRequest(request *models.Request) error {
	existingRequest, err := s.repository.FindByID(request.ID)
	if err != nil {
		return err
	}
	if existingRequest == nil {
		return errors.New("request not found")
	}
	return s.repository.Update(request)
}

func (s *RequestService) GetRequestByID(id uint) (*models.Request, error) {
	return s.repository.FindByID(id)
}

func (s *RequestService) GetAllRequests(limit, offset int) ([]models.Request, int64, error) {
	return s.repository.GetAllRequests(limit, offset)
}

type RequestWithStats struct {
	Request models.Request `json:"request"`
}

func (s *RequestService) GetRequestWithStats(id uint) (*RequestWithStats, error) {
	request, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, errors.New("permintaan tidak ditemukan")
	}
	return &RequestWithStats{
		Request: *request,
	}, nil
}

func (s *RequestService) DeleteRequest(id uint) error {
	request, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if request == nil {
		return errors.New("request not found")
	}
	return s.repository.DeleteByID(id)
}
