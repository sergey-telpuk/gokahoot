package graphql

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
)

func (r *mutationResolver) ActivateGame(ctx context.Context, testUUID string) (*Game, error) {
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)
	testService := r.Di.Container.Get(ContainerNameTestService).(*TestService)
	uuid := guuid.New()

	test, err := testService.FindByUuid(testUUID)
	if err != nil {
		return nil, err
	}

	if err := gameService.CreateNewGame(test.ID, uuid); err != nil {
		return nil, err
	}

	game, err := gameService.FindByCode(uuid.String())

	if err != nil {
		return nil, err
	}

	return mapGame(game)
}
func (r *mutationResolver) DeactivateGameByCODEs(ctx context.Context, codes []string) (*Status, error) {
	service := r.Di.Container.Get(ContainerNameGameService).(*GameService)
	if err := service.DeleteByCODEs(codes...); err != nil {
		return nil, err
	}
	return &Status{Success: true}, nil
}

func (r *mutationResolver) JoinPlayerToGame(ctx context.Context, input InputJoinPlayer) (*Player, error) {
	uuid := guuid.New()
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)

	if status, err := gameService.IsWaitingForJoining(input.GameCode); !status || err != nil {
		return nil, errors.New(fmt.Sprintf("Joing player error: %v", err))
	}

	game, err := gameService.FindByCode(input.GameCode)

	if err != nil {
		return nil, err
	}

	if err := playerService.CreateNewPlayer(uuid, game.ID, input.Name); err != nil {
		return nil, err
	}

	player, err := playerService.FindByUuid(uuid.String())

	if err != nil {
		return nil, err
	}

	if err := broadcastService.BroadcastForWaitForJoiningGame(game.Code, player.UUID); err != nil {
		fmt.Println(errors.New(fmt.Sprintf("Broadcast error: %s", err)))
	}

	return mapPlayer(player)
}
