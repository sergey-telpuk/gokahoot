package services

import (
	"context"
	"errors"
)

func (r *gameResolver) Players(ctx context.Context, game *Game) ([]*Player, error) {
	var rModels []*Player
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)

	mGame, err := gameService.GetGameByCode(game.Code)

	if err != nil {
		return nil, err
	}

	mPlayers, err := playerService.FindPlayersBelongToGame(mGame.ID)

	if err != nil {
		return nil, err
	}

	for _, m := range mPlayers {
		mapped, _ := mapPlayer(m)
		rModels = append(rModels, mapped)
	}

	return rModels, nil
}

func (r *queryResolver) ActivatedGames(ctx context.Context) ([]*Game, error) {
	var rGames []*Game
	service := r.Di.Container.Get(ContainerNameGameService).(*GameService)

	games, err := service.FindAllWitchAreWaitingForJoining()

	if err != nil {
		return nil, err
	}

	for _, m := range games {
		mapped, _ := mapGame(*m)
		rGames = append(rGames, mapped)
	}

	return rGames, nil
}

func (r *queryResolver) ActivatedGameByCode(ctx context.Context, code string) (*Game, error) {
	service := r.Di.Container.Get(ContainerNameGameService).(*GameService)

	game, err := service.GetGameByCode(code)

	if err != nil {
		return nil, err
	}

	ok, err := service.IsWaitingForJoining(game.Code)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("such a game isnt waiting")
	}

	return mapGame(*game)
}
func (r *queryResolver) ReportGameByCode(ctx context.Context, code string) (*ReportGame, error) {

	return nil, nil
}
