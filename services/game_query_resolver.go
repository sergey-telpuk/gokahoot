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

	games, err := service.FindAllWhichAreWaitingForJoining()

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
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)

	game, err := gameService.GetGameByCode(code)

	if err != nil {
		return nil, err
	}

	ok, err := gameService.IsFinishedGame(game.Code)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("such a game isnt finished")
	}

	players, err := playerService.FindPlayersBelongToGame(game.ID)

	if err != nil {
		return nil, err
	}

	game.Players = players

	return mapReport(*game)
}

func (r *queryResolver) ChatMessagesOfGameByCode(ctx context.Context, code string, offset int, limit int, order ChatTimeOrder) ([]*ChatMessage, error) {
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)
	var rChatMessage []*ChatMessage

	game, err := gameService.GetGameByCode(code)

	if err != nil {
		return nil, err
	}

	messages, err := gameService.FindChatMessagesByGameCode(game.Code, offset, limit, order)

	if err != nil {
		return nil, err
	}

	for _, m := range messages {
		mapped, _ := mapChatMessage(*m)
		rChatMessage = append(rChatMessage, mapped)
	}

	return rChatMessage, nil
}
