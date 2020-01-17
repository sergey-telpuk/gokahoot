package graphql

import (
	"context"
	"errors"
	"fmt"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/services"
)

func (r *questionResolver) Answers(ctx context.Context, obj *Question) ([]*Answer, error) {
	var rAnswer []*Answer
	var mAnswer []*models.Answer
	if err := r.Di.Invoke(func(s *services.QuestionService) error {
		var err error
		mAnswer, err = s.FindAnswersBelongToQuestion(obj.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.New(fmt.Sprintf("Provide container was error, error %v", err))
	}

	for _, answer := range mAnswer {
		rAnswer = append(rAnswer, &Answer{
			Text:       answer.Text,
			Sequential: answer.Sequential,
			ImgURL:     answer.ImgURL,
		})
	}

	return rAnswer, nil
}

func (r *testResolver) Questions(ctx context.Context, obj *Test) ([]*Question, error) {
	var rQuestion []*Question
	var mQuestion []*models.Question
	if err := r.Di.Invoke(func(s *services.QuestionService) error {
		var err error
		mQuestion, err = s.FindQuestionBelongToTest(obj.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.New(fmt.Sprintf("Provide container was error, error %v", err))
	}

	for _, q := range mQuestion {
		rQuestion = append(rQuestion, mapQuestion(q))
	}

	return rQuestion, nil
}

func (r *queryResolver) QuestionByID(ctx context.Context, id int) (*Question, error) {
	var mQuestion *models.Question
	if err := r.Di.Invoke(func(s *services.QuestionService) error {
		var err error
		mQuestion, err = s.FindByID(id)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.New(fmt.Sprintf("Provide container was error, error %v", err))
	}

	if mQuestion == nil {
		return nil, errors.New(fmt.Sprintf("Not a such question"))
	}

	return mapQuestion(mQuestion), nil
}

func (r *queryResolver) QuestionByUUID(ctx context.Context, id string) (*Question, error) {
	var mQuestion *models.Question

	if err := r.Di.Invoke(func(s *services.QuestionService) error {
		var err error
		mQuestion, err = s.FindByUuid(id)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.New(fmt.Sprintf("Provide container was error, error %v", err))
	}

	if mQuestion == nil {
		return nil, errors.New(fmt.Sprintf("Not a such question"))
	}

	return mapQuestion(mQuestion), nil
}

func mapQuestion(m *models.Question) *Question {

	return &Question{
		ID:          m.ID,
		UUID:        m.UUID,
		TestID:      m.TestID,
		Text:        m.Text,
		ImgURL:      m.ImgURL,
		RightAnswer: m.RightAnswer,
	}
}
