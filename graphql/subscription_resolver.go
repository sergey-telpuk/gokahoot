package graphql

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"sync"
)

var mutex = &sync.Mutex{}

func (r *subscriptionResolver) OnJoiningPlayerToGame(ctx context.Context, gameCode string, playerUUID string) (<-chan *BroadcastPlayer, error) {
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)
	gametService := r.Di.Container.Get(ContainerNameGameService).(*GameService)

	if status, err := gametService.IsWaitingForJoining(gameCode); !status || err != nil {
		return nil, errors.New(fmt.Sprintf("Joing player error: %v", err))
	}

	uuid := guuid.New()

	if err := broadcastService.AddGame(gameCode); err != nil {
		return nil, err
	}

	player, err := broadcastService.AddPlayerToGame(uuid, gameCode, playerUUID)

	if err != nil {
		return nil, err
	}

	event := make(chan *BroadcastPlayer, 1)
	player.EventWaitForJoining = event

	go func(uuid guuid.UUID, gameCode string) {
		<-ctx.Done()
		mutex.Lock()
		_ = broadcastService.DeletePlayerFromGame(gameCode, uuid)
		mutex.Unlock()
	}(uuid, gameCode)

	return event, nil
}
