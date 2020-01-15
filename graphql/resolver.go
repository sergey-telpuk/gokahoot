package graphql

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
	uuid := guuid.New()
	var test models.TestModel
	if err := r.Di.Invoke(func(testService *services.TestService) {
		testService.CreateNewTest(input.Name, uuid)
		test = testService.FindByUuid(uuid)
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	if test.ID != 0 {
		return &Test{
			Code:      test.Code,
			Name:      test.Name,
			Questions: nil,
		}, nil
	}

	return nil, nil
}

func (r *mutationResolver) CreateNewQuestion(ctx context.Context, input NewQuestion) (*Question, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Tests(ctx context.Context) ([]*Test, error) {
	var mTests []models.TestModel
	var rTests []*Test

	if err := r.Di.Invoke(func(testService *services.TestService) {
		mTests = testService.FindAll()
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}
	if len(mTests) != 0 {
		for _, test := range mTests {
			rTests = append(rTests, &Test{
				Code:      test.Code,
				Name:      test.Name,
				Questions: nil,
			})
		}
		return rTests, nil
	}

	return nil, nil
}
