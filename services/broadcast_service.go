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
		go func(stPlayer *StoragePlayer, player *models.Player, ctx context.Context, cancel context.CancelFunc) {
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
		go func(stPlayer *StoragePlayer, ctx context.Context, cancel context.CancelFunc) {
			select {
			case stPlayer.EventDeletingPlayerFromGame <- &BroadcastPlayer{
				Name:     stPlayer.Name,
				UUID:     stPlayer.UUID,
				GameCode: gameCode,
			}:
			case <-ctx.Done():
				return
			}

			defer cancel()

		}(broadcastPlayer, ctx, cancel)
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
	countQuestions := len(questions)
	commonTime := time.Duration(countQuestions*TimeForAnsweringSec+1) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), commonTime)
	currentTimer := -1
	broadcastTimer := TimeForAnsweringSec + 1
	currentQuestion := questions[0]

	go func() {
		defer cancel()

		for {
			select {
			case <-time.After(1 * time.Second):
				currentTimer++
				broadcastTimer--
				if currentTimer%TimeForAnsweringSec == 0 && currentTimer/TimeForAnsweringSec <= countQuestions {
					currentQuestion = questions[currentTimer/TimeForAnsweringSec]
					broadcastTimer = 5
				}

				s.BroadcastTimerPlayers(broadcastTimer, game.Code, currentQuestion.UUID, gameStatus(game.Status))

			case <-ctx.Done():
				game.Status = models.GameInFinished
				_, _ = s.gameService.Update(&game)
				s.BroadcastTimerPlayers(0, game.Code, currentQuestion.UUID, gameStatus(game.Status))
				return
			}
		}
	}()
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
