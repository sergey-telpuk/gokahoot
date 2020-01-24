package graphql

import (
	"context"
	"github.com/sergey-telpuk/gokahoot/services"
)

func (r *gameResolver) Players(ctx context.Context, game *Game) ([]*Player, error) {
	var rModels []*Player
	playerService := r.Di.Container.Get(services.ContainerNamePlayerService).(*services.PlayerService)
	gameService := r.Di.Container.Get(services.ContainerNameGameService).(*services.GameService)

	mGame, err := gameService.FindByUuid(game.Code)

	if err != nil {
		return nil, err
	}

	mPlayers, err := playerService.FindQuestionBelongToGame(mGame.ID)

	if err != nil {
		return nil, err
	}

	for _, answer := range mPlayers {
		mapped, _ := mapPlayer(answer)
		rModels = append(rModels, mapped)
	}

	return rModels, nil
}
