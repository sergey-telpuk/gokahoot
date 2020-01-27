package graphql

import (
	"context"
	"github.com/sergey-telpuk/gokahoot/services"
)

func (r *gameResolver) Players(ctx context.Context, game *Game) ([]*Player, error) {
	var rModels []*Player
	playerService := r.Di.Container.Get(services.ContainerNamePlayerService).(*services.PlayerService)
	gameService := r.Di.Container.Get(services.ContainerNameGameService).(*services.GameService)

	mGame, err := gameService.FindByCode(game.Code)

	if err != nil {
		return nil, err
	}

	mPlayers, err := playerService.FindQuestionBelongToGame(mGame.ID)

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
	service := r.Di.Container.Get(services.ContainerNameGameService).(*services.GameService)

	games, err := service.FindAll()

	if err != nil {
		return nil, err
	}

	for _, m := range games {
		mapped, _ := mapGame(m)
		rGames = append(rGames, mapped)
	}

	return rGames, nil
}

func (r *queryResolver) ActivatedGameByCode(ctx context.Context, code string) (*Game, error) {
	service := r.Di.Container.Get(services.ContainerNameGameService).(*services.GameService)

	game, err := service.FindByCode(code)

	if err != nil {
		return nil, err
	}

	return mapGame(game)
}
