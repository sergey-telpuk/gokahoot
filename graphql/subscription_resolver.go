package graphql

import (
	"context"
	"fmt"
	"sync"
)

var mutex = &sync.Mutex{}

func (r *subscriptionResolver) OnJoiningPlayerToGame(ctx context.Context, gameCode string, playerUUID string) (<-chan *BroadcastPlayer, error) {
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)
	event := make(chan *BroadcastPlayer, 1)

	if err := broadcastService.AddGame(gameCode); err != nil {
		return nil, err
	}

	player, err := broadcastService.AddPlayerToGame(gameCode, playerUUID)

	if err != nil {
		return nil, err
	}

	player.EventWaitForJoining = event
	fmt.Println(player.EventWaitForJoining)
	go func() {
		<-ctx.Done()
		mutex.Lock()
		fmt.Print("======================")
		_ = broadcastService.DeletePlayerFromGame(gameCode, playerUUID)
		mutex.Unlock()
	}()

	return event, nil
}
