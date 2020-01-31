package services

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
)

const ContainerNameAnswerService = "ContainerNameAnswerService"

type AnswerService struct {
	r *repositories.AnswerRepository
}

func (s AnswerService) FindByID(id int) (*models.Answer, error) {
	return s.r.FindOne("id = ?", id)
}

func (s AnswerService) GetAnswerByID(id int) (*models.Answer, error) {
	m, err := s.FindByID(id)

	if err != nil {
		return nil, err
	}

	if m == nil {
		return nil, errors.New(fmt.Sprintf("Answer model error: %s", "not a such answer"))
	}

	return m, nil
}

func (s AnswerService) UpdateByUUID(m *models.Answer) (*models.Answer, error) {
	return s.r.Update(m)
}

func (s AnswerService) DeleteByUUIDs(id ...string) error {
	return s.r.Delete("uuid IN (?)", id)
}

func (s AnswerService) DeleteByIDs(id ...int) error {
	return s.r.Delete("id IN (?)", id)
}

func InitAnswerService(r *repositories.AnswerRepository) *AnswerService {
	return &AnswerService{
		r: r,
	}
}
