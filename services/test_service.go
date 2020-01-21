package services

import (
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
)

const ContainerNameTestService = "ContainerNameTestService"

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

func (s *TestService) FindByUuid(id string) (*models.Test, error) {
	return s.r.FindOne("uuid = ?", id)
}

func (s *TestService) FindByID(id int) (*models.Test, error) {
	return s.r.FindOne("id = ?", id)
}

func (s *TestService) UpdateByUUID(m *models.Test) (*models.Test, error) {
	return s.r.Update(m)
}

func (s *TestService) DeleteByUUIDs(id ...string) error {
	return s.r.Delete("uuid IN (?)", id)
}

func (s *TestService) DeleteByIDs(id ...int) error {
	return s.r.Delete("id IN (?)", id)
}

func (s *TestService) FindAll() ([]*models.Test, error) {
	return s.r.FindAll()
}

func InitTestService(r *repositories.TestRepository) *TestService {
	return &TestService{
		r: r,
	}
}
