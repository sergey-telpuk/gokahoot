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
	TimeForAnsweringSec           = 15
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
	players, err := s.broadcastRepository.GetPlayersWaitForJoiningGame(gameCode)

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

func (s *BroadcastService) BroadcastForWaitForStartingGame(gameCode string) error {
	players, err := s.broadcastRepository.GetPlayersWaitForStartingGame(gameCode)

	if err != nil {
		return err
	}

	for _, broadcastPlayer := range players {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		go func(_stPlayer *StoragePlayer, gameCode string, _ctx context.Context, _cancel context.CancelFunc) {
			select {
			case _stPlayer.EventStartGame <- &StartGame{
				GameCode: gameCode,
			}:
			case <-_ctx.Done():
				return
			}

			defer _cancel()

		}(broadcastPlayer, gameCode, ctx, cancel)
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

func (s *BroadcastService) BroadcastMessageToChatOFGame(message *models.ChatMessage) error {
	players, err := s.broadcastRepository.GetPlayersForSendingToChatOfGame(message.Game.Code)
	if err != nil {
		return err
	}

	for _, broadcastPlayer := range players {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		go func(_stPlayer *StoragePlayer, _ctx context.Context, _cancel context.CancelFunc) {
			defer _cancel()

			select {
			case _stPlayer.EventChatGame <- &BroadcasChatGame{
				UUID:    message.UUID,
				Message: message.Message,
				Player: &BroadcastPlayer{
					Name:     message.Player.Name,
					UUID:     message.Player.UUID,
					GameCode: message.Player.Game.Code,
				},
				Time: message.CreatedAt.String(),
			}:
			case <-_ctx.Done():
				return
			}
		}(broadcastPlayer, ctx, cancel)
	}

	return nil
}

func (s *BroadcastService) BroadcastIsTypingToChatOFGame(player *models.Player) error {
	players, err := s.broadcastRepository.GetPlayersForChattingGame(player.Game.Code)
	if err != nil {
		return err
	}

	for _, broadcastPlayer := range players {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		go func(_stPlayer *StoragePlayer, _ctx context.Context, _cancel context.CancelFunc) {
			defer _cancel()

			select {
			case _stPlayer.EventIsTypingPlayer <- &BroadcastPlayer{
				UUID:     player.UUID,
				GameCode: player.Game.Code,
				Name:     player.Name,
			}:
			case <-_ctx.Done():
				return
			}
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
	commonTime := time.Duration(len(game.Test.Questions)*TimeForAnsweringSec+TimeForAnsweringSec) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), commonTime)

	go func(_game models.Game, _ctx context.Context, _cancel context.CancelFunc) {
		defer _cancel()

		countQuestions := len(game.Test.Questions)
		if countQuestions == 0 {
			return
		}
		currentTimer := 0
		broadcastTimer := TimeForAnsweringSec
		currentQuestion := game.Test.Questions[0]

		for {
			select {
			case <-time.After(1 * time.Second):
				index := currentTimer / TimeForAnsweringSec
				if currentTimer <= countQuestions*TimeForAnsweringSec && index < countQuestions {
					if currentTimer%TimeForAnsweringSec == 0 {
						currentQuestion = game.Test.Questions[index]
						broadcastTimer = TimeForAnsweringSec
					}

					s.BroadcastTimerPlayers(broadcastTimer, _game, currentQuestion, gameStatus(_game.Status))

					broadcastTimer--
					currentTimer++
					continue
				}

				_cancel()
			case <-_ctx.Done():
				_game.Status = models.GameInFinished
				_, _ = s.gameService.Update(&_game)
				s.BroadcastTimerPlayers(0, _game, currentQuestion, gameStatus(_game.Status))
				return
			}
		}
	}(game, ctx, cancel)
}

func (s *BroadcastService) BroadcastTimerPlayers(timer int, game models.Game, question models.Question, status GameStatus) {
	players, err := s.broadcastRepository.GetPlayersForPlayingGame(game.Code)

	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, player := range players {
		select {
		case player.EventPlayingGame <- &BroadcastPlayingGame{
			CurrentTimeSec:      timer,
			StartTimeSec:        TimeForAnsweringSec,
			GameCode:            game.Code,
			GameStatusEnum:      status,
			CurrentQuestionUUID: question.UUID,
			Answers:             s.buildBroadcastAnswerForChartGame(game, question),
		}:
		case <-ctx.Done():
			return
		}
	}
}

func (s *BroadcastService) buildBroadcastAnswerForChartGame(game models.Game, question models.Question) []*BroadcastAnswerForChartGame {
	answers, _ := s.playerService.FindPlayerAnswersByGameAndQuestion(game, question)
	mapAnswer := map[int][]*BroadcastPlayersForChartGame{}
	var broadcastAnswerForChartGame []*BroadcastAnswerForChartGame

	for _, answer := range answers {
		mapAnswer[answer.AnswerID] = append(mapAnswer[answer.AnswerID], &BroadcastPlayersForChartGame{
			Player: &BroadcastPlayer{
				UUID:     answer.Player.UUID,
				GameCode: answer.Game.Code,
				Name:     answer.Player.Name,
			},
			WasRight: answer.WasRight,
		})
	}

	for answerID, player := range mapAnswer {
		broadcastAnswerForChartGame = append(broadcastAnswerForChartGame, &BroadcastAnswerForChartGame{
			AnswerID: answerID,
			Players:  player,
		})
	}

	return broadcastAnswerForChartGame

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
