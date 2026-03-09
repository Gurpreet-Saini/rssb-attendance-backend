package services

import (
	"attendance-system/backend/models"
	"attendance-system/backend/repositories"
)

type CenterService interface {
	List() ([]models.Center, error)
	Create(name, location string) (*models.Center, error)
	Delete(id uint) error
}

type centerService struct {
	repo repositories.CenterRepository
}

func NewCenterService(repo repositories.CenterRepository) CenterService {
	return &centerService{repo: repo}
}

func (s *centerService) List() ([]models.Center, error) {
	return s.repo.List()
}

func (s *centerService) Create(name, location string) (*models.Center, error) {
	center := &models.Center{Name: name, Location: location}
	if err := s.repo.Create(center); err != nil {
		return nil, err
	}
	return center, nil
}

func (s *centerService) Delete(id uint) error {
	return s.repo.SoftDeleteCascade(id)
}
