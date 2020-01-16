package graphql

//go:generate go run github.com/99designs/gqlgen
import (
	"context"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/services"
	"go.uber.org/dig"
	"log"
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

type questionResolver struct{ *Resolver }

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateNewTest(ctx context.Context, input NewTest) (*Test, error) {
	var test models.Test
	uuid := guuid.New()

	if err := r.Di.Invoke(func(s *services.TestService) {
		s.CreateNewTest(uuid, input.Name, guuid.New())
		test = s.FindByUuid(uuid)
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
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
	var question models.Question
	uuid := guuid.New()
	if err := r.Di.Invoke(func(s *services.QuestionService) {
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
			log.Fatalf("Craeting question was failed, error %v", err)
		}

		question = s.FindByUuid(uuid)

	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
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

	if err := r.Di.Invoke(func(s *services.TestService) {
		mTests = s.FindAll()
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	for _, test := range mTests {
		rTests = append(rTests, &Test{
			ID:        test.ID,
			Code:      test.Code,
			Name:      test.Name,
			Questions: nil,
		})
	}

	return rTests, nil
}

func (r *questionResolver) Answers(ctx context.Context, obj *Question) ([]*Answer, error) {
	var rAnswer []*Answer
	var mAnswer []*models.Answer
	if err := r.Di.Invoke(func(s *services.QuestionService) {
		mAnswer = s.FindAnswersBelongToQuestion(obj.ID)
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
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
