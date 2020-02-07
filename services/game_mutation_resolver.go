package services

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
)

func (r *mutationResolver) ActivateGame(ctx context.Context, testUUID string) (*Game, error) {
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)
	testService := r.Di.Container.Get(ContainerNameTestService).(*TestService)
	uuid := guuid.New()

	test, err := testService.GetTestByUUID(testUUID)
	if err != nil {
		return nil, err
	}

	if err := gameService.CreateNewGame(test.ID, uuid); err != nil {
		return nil, err
	}

	game, err := gameService.GetGameByCode(uuid.String())

	if err != nil {
		return nil, err
	}

	return mapGame(*game)
}
func (r *mutationResolver) DeactivateGameByCODEs(ctx context.Context, codes []string) (*Status, error) {
	service := r.Di.Container.Get(ContainerNameGameService).(*GameService)
	if err := service.DeleteByCODEs(codes...); err != nil {
		return nil, err
	}
	return &Status{Success: true}, nil
}

func (r *mutationResolver) StartGameByCode(ctx context.Context, code string) (*Game, error) {
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)

	game, err := gameService.GetGameByCode(code)

	if err != nil {
		return nil, err
	}

	if status, err := gameService.IsPlayingGame(game.Code); status || err != nil {
		return nil, errors.New(fmt.Sprintf("A statring player messsage: %v or error %v", "a game has already started", err))
	}

	game.Status = models.GameInPlaying

	_, err = gameService.Update(game)

	if err != nil {
		return nil, err
	}

	if err := broadcastService.StartBroadcastGameIsBeingPlayed(game.Code); err != nil {
		fmt.Println(errors.New(fmt.Sprintf("Broadcast error: %s", err)))
	}

	if err := broadcastService.BroadcastForWaitForStartingGame(game.Code); err != nil {
		fmt.Println(errors.New(fmt.Sprintf("Broadcast error: %s", err)))
	}

	return mapGame(*game)
}

func (r *mutationResolver) JoinPlayerToGame(ctx context.Context, input InputJoinPlayer) (*Player, error) {
	uuid := guuid.New()
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)

	if status, err := gameService.IsWaitingForJoining(input.GameCode); !status || err != nil {
		return nil, errors.New(fmt.Sprintf("A joinging player messsage: %v or error %v", "a game isnt acivated", err))
	}

	game, err := gameService.GetGameByCode(input.GameCode)

	if err != nil {
		return nil, err
	}

	if err := playerService.CreateNewPlayer(uuid, game.ID, input.Name); err != nil {
		return nil, err
	}

	player, err := playerService.GetPlayerByUUID(uuid.String())

	if err != nil {
		return nil, err
	}

	if err := broadcastService.BroadcastForWaitForJoiningGame(game.Code, player.UUID); err != nil {
		fmt.Println(errors.New(fmt.Sprintf("Broadcast error: %s", err)))
	}

	return mapPlayer(*player)
}

func (r *mutationResolver) DeletePlayerFromGame(ctx context.Context, gameCode string, playerUUID string) (*Status, error) {
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)

	_, err := playerService.GetPlayerByUUID(playerUUID)

	if err != nil {
		return nil, err
	}

	if err := playerService.DeleteByUUIDs(playerUUID); err != nil {
		return nil, err
	}

	if err := broadcastService.BroadcastForDeletingPlayerGame(gameCode, playerUUID); err != nil {
		fmt.Println(errors.New(fmt.Sprintf("Broadcast error: %s", err)))
	}

	return &Status{Success: true}, nil
}

func (r *mutationResolver) AnswerQuestionByUUID(ctx context.Context, playerUUID string, questionUUID string, answerID int) (*bool, error) {
	questionService := r.Di.Container.Get(ContainerNameQuestionService).(*QuestionService)
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)
	answerService := r.Di.Container.Get(ContainerNameAnswerService).(*AnswerService)
	right := false

	player, err := playerService.GetPlayerByUUID(playerUUID)

	if err != nil {
		return nil, err
	}

	ok, err := gameService.IsPlayingGame(player.Game.Code)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New(fmt.Sprintf("Game: %s", "such a game isnt being played"))
	}

	question, err := questionService.GetQuestionByUUID(questionUUID)

	if err != nil {
		return nil, err
	}

	wasAnswered, err := playerService.FindOnePlayerAnswerByGameAndQuestion(player.Game, *player, *question)

	if err != nil && err.Error() != "PlayerAnswer model error: record not found" {
		return nil, err
	}

	if wasAnswered != nil {
		return nil, errors.New(fmt.Sprintf("Game: %s", "you've aready answerd this question"))
	}

	answer, err := answerService.GetAnswerByID(answerID)

	if err != nil {
		return nil, err
	}

	if question.RightAnswer == answer.Sequential {
		right = true
	}

	if err := playerService.CreateNewPlayerAnswer(player.ID, player.Game.ID, question.ID, answer.ID, right); err != nil {
		return nil, err
	}

	return &right, nil
}

func (r *mutationResolver) SendMessageToChat(ctx context.Context, playerUUID string, message string) (*ChatMessage, error) {
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)
	gameService := r.Di.Container.Get(ContainerNameGameService).(*GameService)
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)
	uuid := guuid.New()

	player, err := playerService.GetPlayerByUUID(playerUUID)

	if err != nil {
		return nil, err
	}

	if err := gameService.CreateNewMessageOfChat(uuid, player.GameID, player.ID, message); err != nil {
		return nil, err
	}

	chatMessage, err := gameService.GetChatMessageByUUID(uuid.String())

	if err != nil {
		return nil, err
	}

	if err := broadcastService.BroadcastMessageToChatOFGame(chatMessage); err != nil {
		fmt.Println(errors.New(fmt.Sprintf("Broadcast error: %s", err)))
	}

	return mapChatMessage(*chatMessage)
}

func (r *mutationResolver) PlayerIsTypingOfGameByUUID(ctx context.Context, playerUUID string) (bool, error) {
	playerService := r.Di.Container.Get(ContainerNamePlayerService).(*PlayerService)
	broadcastService := r.Di.Container.Get(ContainerNameBroadcastService).(*BroadcastService)

	player, err := playerService.GetPlayerByUUID(playerUUID)

	if err != nil {
		return false, err
	}

	if err := broadcastService.BroadcastIsTypingToChatOFGame(player); err != nil {
		fmt.Println(errors.New(fmt.Sprintf("Broadcast error: %s", err)))
	}

	return true, nil
}
