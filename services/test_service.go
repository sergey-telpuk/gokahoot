package services

import (
	"errors"
	"fmt"
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
		Name: name,
	}

	return s.r.Create(model)
}

func (s *TestService) GetTestByUUID(code string) (*models.Test, error) {
	m, err := s.FindByUuid(code)

	if err != nil {
		return nil, err
	}

	if m == nil {
		return nil, errors.New(fmt.Sprintf("Test model error: %s", "not a such test"))
	}

	return m, nil
}

func (s *TestService) FindByUuid(id string) (*models.Test, error) {
	return s.r.FindOne("tests.uuid = ?", id)
}

func (s *TestService) FindByID(id int) (*models.Test, error) {
	return s.r.FindOne("tests.id = ?", id)
}

func (s *TestService) UpdateByUUID(m *models.Test) (*models.Test, error) {
	return s.r.Update(m)
}

func (s *TestService) DeleteByUUIDs(id ...string) error {
	return s.r.Delete("tests.uuid IN (?)", id)
}

func (s *TestService) DeleteByIDs(id ...int) error {
	return s.r.Delete("tests.id IN (?)", id)
}

func (s *TestService) FindAll() ([]*models.Test, error) {
	return s.r.FindAll()
}

func InitTestService(r *repositories.TestRepository) *TestService {
	return &TestService{
		r: r,
	}
}
