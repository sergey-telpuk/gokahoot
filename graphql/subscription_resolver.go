package graphql

import (
	"context"
	"github.com/sergey-telpuk/gokahoot/services"
	"sync"
)

var mutex = &sync.Mutex{}

func (r *subscriptionResolver) OnJoiningPlayerToGame(ctx context.Context, gameCode string, playerUUID string) (<-chan *BroadcastPlayer, error) {
	broadcastService := r.Di.Container.Get(services.ContainerNameBroadcastService).(*services.BroadcastService)
	event := make(chan *BroadcastPlayer, 1)

	if err := broadcastService.AddGame(gameCode); err != nil {
		return nil, err
	}

	_, err := broadcastService.AddPlayerToGame(gameCode, playerUUID)

	if err != nil {
		return nil, err
	}

	//player.EventWaitForJoining = event

	go func() {
		<-ctx.Done()
		mutex.Lock()
		_ = broadcastService.DeletePlayerFromGame(gameCode, playerUUID)
		mutex.Unlock()
	}()

	return event, nil
}
