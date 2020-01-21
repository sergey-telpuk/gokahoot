package services

import (
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
)

const ContainerNameAnswerService = "ContainerNameAnswerService"

type AnswerService struct {
	r *repositories.AnswerRepository
}

func (s *AnswerService) FindByUuid(id string) (*models.Answer, error) {
	return s.r.FindOne("uuid = ?", id)
}

func (s *AnswerService) FindByID(id int) (*models.Answer, error) {
	return s.r.FindOne("id = ?", id)
}

func (s *AnswerService) UpdateByUUID(m *models.Answer) (*models.Answer, error) {
	return s.r.Update(m)
}

func (s *AnswerService) DeleteByUUIDs(id ...string) error {
	return s.r.Delete("uuid IN (?)", id)
}

func (s *AnswerService) DeleteByIDs(id ...int) error {
	return s.r.Delete("id IN (?)", id)
}

func InitAnswerService(r *repositories.AnswerRepository) *AnswerService {
	return &AnswerService{
		r: r,
	}
}
