package graphql

import (
	"errors"
	"fmt"
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

	for _, bPlayer := range game.Players {
		fmt.Println("============", bPlayer, player.Game)
		bPlayer.EventWaitForJoining <- &BroadcastPlayer{
			Name:     player.Name,
			UUID:     player.UUID,
			GameCode: player.Game.Code,
		}
		fmt.Println("============", bPlayer, player.Game)
	}

	return nil
}

func (s *BroadcastService) AddPlayerToGame(gameCode string, playerUUID string) (*StoragePlayer, error) {
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
	if err := s.broadcastRepository.AddPlayerToGame(StoragePlayer{
		UUID:     player.UUID,
		Name:     player.Name,
		GameCode: player.Game.Code,
	}); err != nil {
		return nil, err
	}

	return s.broadcastRepository.GetPlayer(gameCode, playerUUID)
}

func (s *BroadcastService) DeleteGame(gameCode string) error {
	s.broadcastRepository.DeleteGame(gameCode)

	return nil
}

func (s *BroadcastService) DeletePlayerFromGame(gameCode string, playerUUID string) error {
	s.broadcastRepository.DeletePlayerFromGame(gameCode, playerUUID)

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
