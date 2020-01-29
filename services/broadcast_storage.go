package services

import (
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
)

const ContainerNameBroadcastRepository = "ContainerNameBroadcastRepository"

type (
	BroadcastRepository struct {
		Games map[string]*StorageGame
	}

	StorageGame struct {
		Code    string
		Players map[string]*StoragePlayer
	}

	StoragePlayer struct {
		Name                        string
		UUID                        string
		GameCode                    string
		EventWaitForJoining         chan *BroadcastPlayer
		EventDeletingPlayerFromGame chan *BroadcastPlayer
		EventPlayingGame            chan *BroadcastPlayingGame
		EventStartGame              chan *StartGame
	}
)

func (r *BroadcastRepository) AddGame(gameCode string) {
	_, ok := r.Games[gameCode]

	if !ok {
		r.Games[gameCode] = &StorageGame{
			Code:    gameCode,
			Players: map[string]*StoragePlayer{},
		}
	}
}

func (r *BroadcastRepository) GetGame(gameCode string) (*StorageGame, error) {
	if !r.HasGame(gameCode) {
		return nil, errorBroadcastRepository(errors.New(fmt.Sprintf("not such a game(%s)", gameCode)))
	}

	return r.Games[gameCode], nil
}

func (r *BroadcastRepository) GetPlayersForPlayingGame(gameCode string) ([]*StoragePlayer, error) {
	var players []*StoragePlayer
	game, err := r.GetGame(gameCode)

	if err != nil {
		return nil, err
	}

	for _, _player := range game.Players {
		if _player.EventPlayingGame == nil {
			continue
		}
		players = append(players, _player)
	}

	return players, nil
}

func (r *BroadcastRepository) GetPlayersWaitForJoiningGame(gameCode string) ([]*StoragePlayer, error) {
	var players []*StoragePlayer
	game, err := r.GetGame(gameCode)

	if err != nil {
		return nil, err
	}

	for _, _player := range game.Players {
		if _player.EventWaitForJoining == nil {
			continue
		}
		players = append(players, _player)
	}

	return players, nil
}

func (r *BroadcastRepository) GetPlayersWaitForStartingGame(gameCode string) ([]*StoragePlayer, error) {
	var players []*StoragePlayer
	game, err := r.GetGame(gameCode)

	if err != nil {
		return nil, err
	}

	for _, _player := range game.Players {
		if _player.EventStartGame == nil {
			continue
		}
		players = append(players, _player)
	}

	return players, nil
}

func (r *BroadcastRepository) GetPlayersForDeletingGame(gameCode string) ([]*StoragePlayer, error) {
	var players []*StoragePlayer
	game, err := r.GetGame(gameCode)

	if err != nil {
		return nil, err
	}

	for _, _player := range game.Players {
		if _player.EventDeletingPlayerFromGame == nil {
			continue
		}
		players = append(players, _player)
	}

	return players, nil
}

func (r *BroadcastRepository) DeleteGame(gameCode string) {
	_, ok := r.Games[gameCode]

	if ok {
		delete(r.Games, gameCode)
	}
}

func (r *BroadcastRepository) HasGame(gameCode string) bool {
	_, ok := r.Games[gameCode]

	return ok
}

func (r *BroadcastRepository) AddPlayerToGame(uuid guuid.UUID, player *StoragePlayer) error {
	game, err := r.GetGame(player.GameCode)

	if err != nil {
		return err
	}

	game.Players[uuid.String()] = player

	return nil
}

func (r *BroadcastRepository) DeletePlayerFromGame(gameCode string, uuid guuid.UUID) {
	_, ok := r.Games[gameCode].Players[uuid.String()]

	if ok {
		delete(r.Games[gameCode].Players, uuid.String())
	}

}

func InitBroadcastRepository() *BroadcastRepository {
	return &BroadcastRepository{Games: map[string]*StorageGame{}}
}

func errorBroadcastRepository(err error) error {
	return errors.New(fmt.Sprintf("BroadcastStorage error: %s", err))
}
