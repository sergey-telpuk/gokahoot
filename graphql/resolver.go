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
		s.CreateNewQuestion(
			uuid,
			input.TestID,
			input.Text,
			input.ImgURL,
			input.RightAnswer,
		)

		question = s.FindByUuid(uuid)

	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	if question.ID != 0 {
		return &Question{
			ID:          question.ID,
			TestID:      question.TestID,
			Text:        question.Text,
			ImgURL:      question.ImgURL,
			RightAnswer: question.RightAnswer,
			Answers:     nil,
		}, nil
	}

	return nil, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Tests(ctx context.Context) ([]*Test, error) {
	var mTests []models.Test
	var rTests []*Test

	if err := r.Di.Invoke(func(testService *services.TestService) {
		mTests = testService.FindAll()
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}
	if len(mTests) != 0 {
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

	return nil, nil
}
