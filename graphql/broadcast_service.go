package graphql

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"time"
)

const ContainerNameBroadcastService = "ContainerNameBroadcastService"

type BroadcastService struct {
	broadcastRepository *BroadcastRepository
	gameService         *GameService
	playerService       *PlayerService
}

func (s *BroadcastService) AddGame(gameCode string) error {
	game, err := s.gameService.FindByCode(gameCode)

	if err != nil {
		return errorBroadcastService(err)
	}

	if game == nil {
		return errorBroadcastService(errors.New("not such a game"))
	}

	s.broadcastRepository.AddGame(gameCode)

	return nil
}

func (s *BroadcastService) BroadcastForWaitForJoiningGame(gameCode string, playerUUID string) error {
	game, err := s.broadcastRepository.GetGame(gameCode)

	if err != nil {
		return err
	}

	player, err := s.playerService.FindByUuid(playerUUID)

	if err != nil {
		return err
	}

	for _, broadcastPlayer := range game.Players {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		go func(stPlayer *StoragePlayer, player *models.Player, ctx context.Context) {
			select {
			case stPlayer.EventWaitForJoining <- &BroadcastPlayer{
				Name:     player.Name,
				UUID:     player.UUID,
				GameCode: player.Game.Code,
			}:
			case <-ctx.Done():
				return
			}

			defer cancel()

		}(broadcastPlayer, player, ctx)
	}

	return nil
}

func (s *BroadcastService) AddPlayerToGame(uuid guuid.UUID, gameCode string, playerUUID string) (*StoragePlayer, error) {
	ok := s.broadcastRepository.HasGame(gameCode)

	if !ok {
		return nil, errorBroadcastService(errors.New("not such a game"))
	}

	player, err := s.playerService.FindByUuid(playerUUID)

	if err != nil {
		return nil, errorBroadcastService(err)
	}

	if player == nil {
		return nil, errorBroadcastService(errors.New("not such a player"))
	}

	newPlayer := &StoragePlayer{
		UUID:     player.UUID,
		Name:     player.Name,
		GameCode: player.Game.Code,
	}

	if err := s.broadcastRepository.AddPlayerToGame(uuid, newPlayer); err != nil {
		return nil, err
	}

	return newPlayer, nil
}

func (s *BroadcastService) DeleteGame(gameCode string) error {
	s.broadcastRepository.DeleteGame(gameCode)

	return nil
}

func (s *BroadcastService) DeletePlayerFromGame(gameCode string, uuid guuid.UUID) error {
	s.broadcastRepository.DeletePlayerFromGame(gameCode, uuid)

	return nil
}

func InitBroadcastService(
	broadcast *BroadcastRepository,
	gameService *GameService,
	playerService *PlayerService,

) *BroadcastService {
	return &BroadcastService{
		broadcastRepository: broadcast,
		gameService:         gameService,
		playerService:       playerService,
	}
}

func errorBroadcastService(err error) error {
	return errors.New(fmt.Sprintf("Broadcast error: %s", err))
}
