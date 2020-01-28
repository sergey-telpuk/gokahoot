package services

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"time"
)

const (
	ContainerNameBroadcastService = "ContainerNameBroadcastService"
	TimeForAnsweringSec           = 5
)

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

func (s *BroadcastService) StartBroadcastGameIsBeingPlayed(gameCode string) error {
	game, err := s.gameService.GetGameByCode(gameCode)

	if err != nil {
		return err
	}

	game, err = s.gameService.AddRelationsQuestionsAndPlayers(game)

	if err != nil {
		return err
	}

	go s.PlayGame(*game)

	return nil
}

func (s *BroadcastService) BroadcastForWaitForJoiningGame(gameCode string, playerUUID string) error {
	players, err := s.broadcastRepository.GetPlayersForPlayingGame(gameCode)

	if err != nil {
		return err
	}

	player, err := s.playerService.FindByUuid(playerUUID)

	if err != nil {
		return err
	}

	for _, broadcastPlayer := range players {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		go func(_stPlayer *StoragePlayer, _player *models.Player, _ctx context.Context, _cancel context.CancelFunc) {
			select {
			case _stPlayer.EventWaitForJoining <- &BroadcastPlayer{
				Name:     _player.Name,
				UUID:     _player.UUID,
				GameCode: _player.Game.Code,
			}:
			case <-_ctx.Done():
				return
			}

			defer _cancel()

		}(broadcastPlayer, player, ctx, cancel)
	}

	return nil
}

func (s *BroadcastService) BroadcastForDeletingPlayerGame(gameCode string, playerUUID string) error {
	players, err := s.broadcastRepository.GetPlayersForDeletingGame(gameCode)

	if err != nil {
		return err
	}

	for _, broadcastPlayer := range players {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		go func(_stPlayer *StoragePlayer, _gameCode string, _playerUUID string, _ctx context.Context, _cancel context.CancelFunc) {
			defer _cancel()

			select {
			case _stPlayer.EventDeletingPlayerFromGame <- &BroadcastPlayer{
				Name:     "",
				UUID:     _playerUUID,
				GameCode: _gameCode,
			}:
			case <-_ctx.Done():
				return
			}
		}(broadcastPlayer, gameCode, playerUUID, ctx, cancel)
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

func (s *BroadcastService) PlayGame(game models.Game) {
	questions := game.Test.Questions
	commonTime := time.Duration(len(questions)*TimeForAnsweringSec+1) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), commonTime)

	go func(_game models.Game, _questions []models.Question, _ctx context.Context, _cancel context.CancelFunc) {
		defer _cancel()
		countQuestions := len(_questions)
		currentTimer := 0
		broadcastTimer := TimeForAnsweringSec
		currentQuestion := _questions[0]

		for {
			select {
			case <-time.After(1 * time.Second):
				if currentTimer%TimeForAnsweringSec == 0 && currentTimer/TimeForAnsweringSec <= countQuestions {
					currentQuestion = _questions[currentTimer/TimeForAnsweringSec]
					broadcastTimer = TimeForAnsweringSec
				}

				s.BroadcastTimerPlayers(broadcastTimer, _game.Code, currentQuestion.UUID, gameStatus(_game.Status))

				broadcastTimer--
				currentTimer++
			case <-_ctx.Done():
				_game.Status = models.GameInFinished
				_, _ = s.gameService.Update(&_game)
				s.BroadcastTimerPlayers(0, _game.Code, currentQuestion.UUID, gameStatus(_game.Status))
				return
			}
		}
	}(game, questions, ctx, cancel)
}

func (s *BroadcastService) BroadcastTimerPlayers(timer int, gameCode string, questionUUID string, status GameStatus) {
	players, err := s.broadcastRepository.GetPlayersForPlayingGame(gameCode)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err != nil {
		return
	}

	for _, player := range players {
		select {
		case player.EventPlayingGame <- &BroadcastPlayingGame{
			Timer:               timer,
			GameCode:            gameCode,
			GameStatusEnum:      status,
			CurrentQuestionUUID: questionUUID,
		}:
		case <-ctx.Done():
			return
		}
	}
}

func gameStatus(status int) GameStatus {
	switch status {
	case models.GameInPlaying:
		return GameStatusPlaying
	case models.GameInFinished:
		return GameStatusFinished
	default:
		return GameStatusWaitForPlayers
	}
}

func errorBroadcastService(err error) error {
	return errors.New(fmt.Sprintf("Broadcast error: %s", err))
}