package graphql

import (
	"context"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/services"
)

func (r *mutationResolver) ActivateGame(ctx context.Context, testUUID string) (*Game, error) {
	gameService := r.Di.Container.Get(services.ContainerNameGameService).(*services.GameService)
	testService := r.Di.Container.Get(services.ContainerNameTestService).(*services.TestService)
	uuid := guuid.New()

	test, err := testService.FindByUuid(testUUID)
	if err != nil {
		return nil, err
	}

	if err := gameService.CreateNewGame(test.ID, uuid); err != nil {
		return nil, err
	}

	game, err := gameService.FindByUuid(uuid.String())

	if err != nil {
		return nil, err
	}

	return mapGame(game)
}
func (r *mutationResolver) DeactivateGameByCODEs(ctx context.Context, codes []string) (*Status, error) {
	service := r.Di.Container.Get(services.ContainerNameGameService).(*services.GameService)
	if err := service.DeleteByUUIDs(codes...); err != nil {
		return nil, err
	}
	return &Status{Message: true}, nil
}
