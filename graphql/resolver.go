package graphql

//go:generate go run github.com/99designs/gqlgen
import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/services"
	"go.uber.org/dig"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	Di *dig.Container
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Question() QuestionResolver {
	return &questionResolver{r}
}

func (r *Resolver) Test() TestResolver {
	return &testResolver{r}
}

type questionResolver struct{ *Resolver }

type testResolver struct{ *Resolver }

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateNewTest(ctx context.Context, input NewTest) (*Test, error) {
	var test *models.Test
	uuid := guuid.New()

	if err := r.Di.Invoke(func(s *services.TestService) error {
		if err := s.CreateNewTest(uuid, input.Name, guuid.New()); err != nil {
			return err
		}
		var err error
		test, err = s.FindByUuid(uuid)

		if err != nil {
			return err
		}

		return nil

	}); err != nil {
		return nil, errors.New(fmt.Sprintf("Provide container was error, error %v", err))
	}

	if test.ID != 0 {
		return &Test{
			ID:        test.ID,
			UUID:      test.UUID,
			Code:      test.Code,
			Name:      test.Name,
			Questions: nil,
		}, nil
	}

	return nil, nil
}

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
		question, err = s.FindByUuid(uuid)

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

type queryResolver struct{ *Resolver }

func (r *queryResolver) Tests(ctx context.Context) ([]*Test, error) {
	var mTests []models.Test
	var rTests []*Test

	if err := r.Di.Invoke(func(s *services.TestService) error {
		var err error
		mTests, err = s.FindAll()
		if err != nil {
			return err
		}

		return nil

	}); err != nil {
		return nil, errors.New(fmt.Sprintf("Provide container was error, error %v", err))
	}

	for _, test := range mTests {
		rTests = append(rTests, &Test{
			ID:   test.ID,
			Code: test.Code,
			Name: test.Name,
		})
	}

	return rTests, nil
}

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
		rQuestion = append(rQuestion, &Question{
			ID:          q.ID,
			UUID:        q.UUID,
			TestID:      q.TestID,
			Text:        q.Text,
			ImgURL:      q.ImgURL,
			RightAnswer: q.RightAnswer,
		})
	}

	return rQuestion, nil
}
