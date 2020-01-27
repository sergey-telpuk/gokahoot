package graphql

import (
	"errors"
	"fmt"
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
		Name                string
		UUID                string
		GameCode            string
		EventWaitForJoining chan *BroadcastPlayer
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

func (r *BroadcastRepository) GetPlayer(gameCode string, playerUUID string) (*StoragePlayer, error) {

	game, err := r.GetGame(gameCode)

	if err != nil {
		return nil, err
	}

	player, ok := game.Players[playerUUID]

	if !ok {
		return nil, errorBroadcastRepository(errors.New(fmt.Sprintf("not such a player(%s)", gameCode)))
	}

	return player, nil
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

func (r *BroadcastRepository) AddPlayerToGame(player StoragePlayer) error {
	game, err := r.GetGame(player.GameCode)

	if err != nil {
		return err
	}

	_, ok := game.Players[player.UUID]

	if !ok {
		game.Players[player.UUID] = &player
	}

	return nil
}

func (r *BroadcastRepository) DeletePlayerFromGame(gameCode string, playerUUID string) {
	_, ok := r.Games[gameCode].Players[playerUUID]

	if ok {
		delete(r.Games[gameCode].Players, playerUUID)
	}

}

func InitBroadcastRepository() *BroadcastRepository {
	return &BroadcastRepository{Games: map[string]*StorageGame{}}
}

func errorBroadcastRepository(err error) error {
	return errors.New(fmt.Sprintf("BroadcastStorage error: %s", err))
}
