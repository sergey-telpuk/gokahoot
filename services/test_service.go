package services

import (
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
)

type TestService struct {
	r *repositories.TestRepository
}

func (s *TestService) CreateNewTest(uuid guuid.UUID, name string, code guuid.UUID) error {
	model := &models.Test{
		UUID: uuid.String(),
		Code: code.String(),
		Name: name,
	}

	return s.r.Create(model)
}

func (s *TestService) FindByUuid(id guuid.UUID) (*models.Test, error) {
	return s.r.FindOne("uuid = ?", id.String())
}

func (s *TestService) FindAll() ([]models.Test, error) {
	return s.r.FindAll()
}

func InitTestService(r *repositories.TestRepository) *TestService {
	return &TestService{
		r: r,
	}
}
