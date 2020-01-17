package graphql

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/services"
)

func (r *mutationResolver) CreateNewQuestion(ctx context.Context, input NewQuestion) (*Question, error) {
	var question *models.Question
	uuid := guuid.New()
	if err := r.Di.Invoke(func(s *services.QuestionService) error {
		var answers []*models.Answer
		for _, answer := range input.Answers {
			answers = append(answers, &models.Answer{
				Text:   answer.Text,
				ImgURL: answer.ImgURL,
			})
		}
		if err := s.CreateNewQuestion(
			uuid,
			input.TestID,
			input.Text,
			input.ImgURL,
			input.RightAnswer,
			answers,
		); err != nil {
			return errors.New(fmt.Sprintf("Craeting question was failed, error %v", err))
		}
		var err error
		question, err = s.FindByUuid(uuid.String())

		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.New(fmt.Sprintf("Provide container was error, error %v", err))
	}

	if question.ID == 0 {
		return nil, nil
	}

	return &Question{
		ID:          question.ID,
		TestID:      question.TestID,
		UUID:        question.UUID,
		Text:        question.Text,
		ImgURL:      question.ImgURL,
		RightAnswer: question.RightAnswer,
	}, nil
}
