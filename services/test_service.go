package services

import (
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
)

type TestService struct {
	repository *repositories.TestRepository
}

func (s *TestService) CreateNewTest(name string, id guuid.UUID) {
	model := &models.TestModel{
		Code: id.String(),
		Name: name,
	}

	s.repository.Create(model)
}
func (s *TestService) FindByUuid(id guuid.UUID) models.TestModel {
	return s.repository.FindOne("code = ?", id)
}
func (s *TestService) FindAll() []models.TestModel {
	return s.repository.FindAll()
}

func InitTestService(repository *repositories.TestRepository) *TestService {
	return &TestService{
		repository: repository,
	}
}
