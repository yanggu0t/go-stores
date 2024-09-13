package services

import (
	"errors"

	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"gorm.io/gorm"
)

type ProjectService struct {
	DB *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{DB: db}
}

func (s *ProjectService) GetAllProjects() ([]models.Project, error) {
	var projects []models.Project
	if err := s.DB.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *ProjectService) CreateProject(project *models.Project) error {
	return s.DB.Create(project).Error
}

func (s *ProjectService) GetProjectByID(id string) (*models.Project, error) {
	var project models.Project
	if err := s.DB.First(&project, "projectId = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}
	return &project, nil
}

func (s *ProjectService) UpdateProject(id string, project *models.Project) error {
	return s.DB.Model(&models.Project{}).Where("projectId = ?", id).Updates(project).Error
}

func (s *ProjectService) DeleteProject(id string) error {
	return s.DB.Delete(&models.Project{}, "projectId = ?", id).Error
}
