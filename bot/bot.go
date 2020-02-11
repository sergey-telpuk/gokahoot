package bot

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/services"
	"math/rand"
	"syreclabs.com/go/faker"
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
	playerService := b.di.Container.Get(services.ContainerNamePlayerService).(*services.PlayerService)
	broadcastService := b.di.Container.Get(services.ContainerNameBroadcastService).(*services.BroadcastService)

	waitFoJoining := func() chan []*models.Game {
		ch := make(chan []*models.Game)
		go func() {
			for {
				select {
				case <-time.After(1 * time.Second):
					games, _ := gameService.FindAllWhichAreWaitingForJoining()

					if len(games) > 0 {
						ch <- games
					}
				}
			}
		}()

		return ch
	}

	playingGames := func() chan []*models.Game {
		ch := make(chan []*models.Game)

		go func() {
			for {
				select {
				case <-time.After(1 * time.Second):
					games, _ := gameService.FindAllWhichArePlaying()

					if len(games) > 0 {
						ch <- games
					}
				}
			}
		}()

		return ch
	}

	go func() {
		for games := range waitFoJoining() {
			for _, game := range games {
				ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

				go func(game *models.Game, ctx context.Context) {
					for {
						select {
						case <-time.After(1 * time.Second):
							go b.joinPlayer(game.Code, faker.Name().Prefix()+" "+faker.Name().Name()+" "+faker.Name().FirstName()+" "+faker.Name().LastName())
						case <-ctx.Done():
							return
						}
					}

				}(game, ctx)
			}

		}
	}()

	go func() {
		for games := range playingGames() {
			for _, game := range games {
				questions := game.Test.Questions
				players, _ := playerService.FindPlayersBelongToGame(game.ID)

				for _, question := range questions {
					answers := question.Answers
					for _, player := range players {
						go func(player models.Player) {
							time.Sleep(50 * time.Microsecond)
							answer := randomAnswer(answers)
							right := false
							if question.RightAnswer == answer.Sequential {
								right = true
							}

							_ = playerService.CreateNewPlayerAnswer(player.ID, player.Game.ID, question.ID, answer.ID, right)
							uuid := guuid.New()
							_ = gameService.CreateNewMessageOfChat(uuid, player.GameID, player.ID, faker.Lorem().Sentence(10))
							chatMessage, _ := gameService.GetChatMessageByUUID(uuid.String())

							if err := broadcastService.BroadcastMessageToChatOFGame(chatMessage); err != nil {
								fmt.Println(errors.New(fmt.Sprintf("Broadcast error: %s", err)))
							}
						}(player)

					}
					time.Sleep(15 * time.Second)
				}
			}

		}
	}()

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

	_ = gameService.CreateNewMessageOfChat(uuid, player.GameID, player.ID, faker.Lorem().Sentence(10))
	chatMessage, _ := gameService.GetChatMessageByUUID(uuid.String())

	if err := broadcastService.BroadcastMessageToChatOFGame(chatMessage); err != nil {
		fmt.Println(errors.New(fmt.Sprintf("Broadcast error: %s", err)))
	}

}

func randomAnswer(answers []models.Answer) models.Answer {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := len(answers)
	n := rand.Intn(max-min) + min

	return answers[n]
}

func (b Bot) joinToGame() {

}

func (b Bot) writeMassageToGame() {

}
