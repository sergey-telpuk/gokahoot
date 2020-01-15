package services

import (
	"github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
)

type QuestionService struct {
	r *repositories.QuestionRepository
}

func (s *QuestionService) CreateNewQuestion(uuid uuid.UUID, testID int, text string, imgURL *string, rightAnswer int) {
	model := &models.Question{
		UUID:        uuid.String(),
		TestID:      testID,
		Text:        text,
		ImgURL:      imgURL,
		RightAnswer: rightAnswer,
	}

	s.r.Create(model)
}

func (s *QuestionService) FindByUuid(id uuid.UUID) models.Question {
	return s.r.FindOne("uuid = ?", id.String())
}

func (s *QuestionService) FindAll() []models.Question {
	return s.r.FindAll()
}

func InitQuestionService(r *repositories.QuestionRepository) *QuestionService {
	return &QuestionService{
		r: r,
	}
}
