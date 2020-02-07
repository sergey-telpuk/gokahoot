package bot

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/services"
	"time"
)

type Bot struct {
	di *services.DI
}

func Init(di *services.DI) *Bot {
	return &Bot{di: di}
}

func (b Bot) Run() {

	b.tryToFindGameForWaitingForJoiningPlayers()
}

func (b Bot) tryToFindGameForWaitingForJoiningPlayers() {
	gameService := b.di.Container.Get(services.ContainerNameGameService).(*services.GameService)
	chGames := make(chan []*models.Game)

	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				games, _ := gameService.FindAllWhichAreWaitingForJoining()

				if len(games) > 1 {
					chGames <- games
					return
				}
			}
		}
	}()

	for _, game := range <-chGames {
		ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
		go func(game *models.Game, ctx context.Context) {
			for {
				select {
				case <-time.After(1 * time.Second):
					go b.joinPlayer(game.Code, "Bot_Sergey")
				case <-ctx.Done():
					return
				}
			}

		}(game, ctx)
	}
}

func (b Bot) joinPlayer(gameCode string, name string) {
	uuid := guuid.New()
	gameService := b.di.Container.Get(services.ContainerNameGameService).(*services.GameService)
	playerService := b.di.Container.Get(services.ContainerNamePlayerService).(*services.PlayerService)
	broadcastService := b.di.Container.Get(services.ContainerNameBroadcastService).(*services.BroadcastService)

	game, _ := gameService.GetGameByCode(gameCode)

	_ = playerService.CreateNewPlayer(uuid, game.ID, name)
	player, _ := playerService.GetPlayerByUUID(uuid.String())

	if err := broadcastService.BroadcastForWaitForJoiningGame(game.Code, player.UUID); err != nil {
		fmt.Println(errors.New(fmt.Sprintf("Broadcast error: %s", err)))
	}

	_ = gameService.CreateNewMessageOfChat(uuid, player.GameID, player.ID, "Hello from Bot")
	chatMessage, _ := gameService.GetChatMessageByUUID(uuid.String())

	if err := broadcastService.BroadcastMessageToChatOFGame(chatMessage); err != nil {
		fmt.Println(errors.New(fmt.Sprintf("Broadcast error: %s", err)))
	}

}

func (b Bot) joinToGame() {

}

func (b Bot) writeMassageToGame() {

}
