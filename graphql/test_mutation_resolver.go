package graphql

import (
	"context"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/services"
)

func (r *mutationResolver) CreateNewTest(ctx context.Context, input NewTest) (*Test, error) {
	uuid := guuid.New()
	service := r.Di.Container.Get(services.ContainerNameTestService).(*services.TestService)

	if err := service.CreateNewTest(uuid, input.Name, guuid.New()); err != nil {
		return nil, err
	}

	test, err := service.FindByUuid(uuid.String())

	if err != nil {
		return nil, err
	}

	return mapTest(test)
}
