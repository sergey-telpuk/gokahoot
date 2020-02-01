package services

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"sync"
)

var mutex = &sync.Mutex{}

func (r *subscriptionResolver) OnWaitForJoiningPlayerToGame(ctx context.Context, gameCode string, playerUUID string) (<-chan *BroadcastPlayer, error) {
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)

	mPlayer, err := playerService.GetPlayerByUUID(playerUUID)

	if err != nil {
		return nil, err
	}

	if mPlayer.Game.Code != gameCode {
		return nil, errors.New(fmt.Sprintf("A joinging player messsage: %v or error %v", "a palyer isnt belonged to game", err))
	}

	if status, err := gameService.IsWaitingForJoining(gameCode); !status || err != nil {
		return nil, errors.New(fmt.Sprintf("A joinging player messsage: %v or error %v", "a game isnt acivated", err))
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

func (r *subscriptionResolver) OnPlayingGame(ctx context.Context, gameCode string, playerUUID string) (<-chan *BroadcastPlayingGame, error) {
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)

	mPlayer, err := playerService.GetPlayerByUUID(playerUUID)

	if err != nil {
		return nil, err
	}

	if mPlayer.Game.Code != gameCode {
		return nil, errors.New(fmt.Sprintf("A playing game messsage: %v or error %v", "a palyer isnt belonged to game", err))
	}

	if status, err := gameService.IsPlayingGame(gameCode); !status || err != nil {
		return nil, errors.New(fmt.Sprintf("A joinging player messsage: %v or error %v", "a game isnt acivated", err))
	}

	uuid := guuid.New()

	if err := broadcastService.AddGame(gameCode); err != nil {
		return nil, err
	}

	player, err := broadcastService.AddPlayerToGame(uuid, gameCode, playerUUID)

	if err != nil {
		return nil, err
	}

	event := make(chan *BroadcastPlayingGame, 1)
	player.EventPlayingGame = event

	go func(uuid guuid.UUID, gameCode string) {
		<-ctx.Done()
		mutex.Lock()
		_ = broadcastService.DeletePlayerFromGame(gameCode, uuid)
		mutex.Unlock()
	}(uuid, gameCode)

	return event, nil
}

func (r *subscriptionResolver) OnDeletePlayerFromGame(ctx context.Context, gameCode string, playerUUID string) (<-chan *BroadcastPlayer, error) {
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)

	mPlayer, err := playerService.GetPlayerByUUID(playerUUID)

	if err != nil {
		return nil, err
	}

	if mPlayer.Game.Code != gameCode {
		return nil, errors.New(fmt.Sprintf("A playing game messsage: %v or error %v", "a palyer isnt belonged to game", err))
	}

	if status, err := gameService.IsPlayingGame(gameCode); status || err != nil {
		return nil, errors.New(fmt.Sprintf("A joinging player messsage: %v or error %v", "a game isnt acivated", err))
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
	player.EventDeletingPlayerFromGame = event

	go func(uuid guuid.UUID, gameCode string) {
		<-ctx.Done()
		mutex.Lock()
		_ = broadcastService.DeletePlayerFromGame(gameCode, uuid)
		mutex.Unlock()
	}(uuid, gameCode)

	return event, nil
}

func (r *subscriptionResolver) OnWaitForStartingGame(ctx context.Context, gameCode string, playerUUID string) (<-chan *StartGame, error) {
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)

	mPlayer, err := playerService.GetPlayerByUUID(playerUUID)

	if err != nil {
		return nil, err
	}

	if mPlayer.Game.Code != gameCode {
		return nil, errors.New(fmt.Sprintf("A playing game messsage: %v or error %v", "a palyer isnt belonged to game", err))
	}

	uuid := guuid.New()

	if err := broadcastService.AddGame(gameCode); err != nil {
		return nil, err
	}

	player, err := broadcastService.AddPlayerToGame(uuid, gameCode, playerUUID)

	if err != nil {
		return nil, err
	}

	event := make(chan *StartGame, 1)
	player.EventStartGame = event

	go func(uuid guuid.UUID, gameCode string) {
		<-ctx.Done()
		mutex.Lock()
		_ = broadcastService.DeletePlayerFromGame(gameCode, uuid)
		mutex.Unlock()
	}(uuid, gameCode)

	return event, nil
}

func (r *subscriptionResolver) OnChatGame(ctx context.Context, gameCode string, playerUUID string) (<-chan *BroadcasChatGame, error) {
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)

	mPlayer, err := playerService.GetPlayerByUUID(playerUUID)

	if err != nil {
		return nil, err
	}

	if mPlayer.Game.Code != gameCode {
		return nil, errors.New(fmt.Sprintf("A playing game messsage: %v or error %v", "a palyer isnt belonged to game", err))
	}

	uuid := guuid.New()

	if err := broadcastService.AddGame(gameCode); err != nil {
		return nil, err
	}

	player, err := broadcastService.AddPlayerToGame(uuid, gameCode, playerUUID)

	if err != nil {
		return nil, err
	}

	event := make(chan *BroadcasChatGame, 1)
	player.EventChatGame = event

	go func(uuid guuid.UUID, gameCode string) {
		<-ctx.Done()
		mutex.Lock()
		_ = broadcastService.DeletePlayerFromGame(gameCode, uuid)
		mutex.Unlock()
	}(uuid, gameCode)

	return event, nil
}
