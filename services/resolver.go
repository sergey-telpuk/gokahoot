package services

//go:generate go run github.com/99designs/gqlgen
import (
	"github.com/sergey-telpuk/gokahoot/models"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	Di *DI
}

type questionResolver struct{ *Resolver }

type testResolver struct{ *Resolver }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

type gameResolver struct{ *Resolver }

type subscriptionResolver struct{ *Resolver }

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Question() QuestionResolver {
	return &questionResolver{r}
}

func (r *Resolver) Game() GameResolver {
	return &gameResolver{r}
}

func (r *Resolver) Test() TestResolver {
	return &testResolver{r}
}

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

func mapQuestion(m models.Question) (*Question, error) {

	return &Question{
		ID:          m.ID,
		UUID:        m.UUID,
		TestID:      m.TestID,
		Text:        m.Text,
		ImgURL:      m.ImgURL,
		RightAnswer: m.RightAnswer,
	}, nil
}

func mapAnswer(m models.Answer) (*Answer, error) {

	return &Answer{
		ID:         m.ID,
		Text:       m.Text,
		Sequential: m.Sequential,
		ImgURL:     m.ImgURL,
	}, nil
}

func mapTest(m models.Test) (*Test, error) {

	return &Test{
		ID:   m.ID,
		UUID: m.UUID,
		Name: m.Name,
	}, nil

}

func mapGame(m models.Game) (*Game, error) {
	test, _ := mapTest(m.Test)

	return &Game{
		Code:   m.Code,
		Status: m.Status,
		Test:   test,
	}, nil

}
func mapPlayer(m models.Player) (*Player, error) {
	game, _ := mapGame(m.Game)

	return &Player{
		UUID: m.UUID,
		Name: m.Name,
		Game: game,
	}, nil

}

func mapReport(m models.Game) (*ReportGame, error) {
	game, _ := mapGame(m)
	var players []*ReportPlayer

	for _, _p := range m.Players {
		var ra []*ReportAnswer

		_mapPl, _ := mapPlayer(_p)
		for _, _pa := range _p.PlayerAnswers {
			_mapAns, _ := mapAnswer(_pa.Answer)
			ra = append(ra, &ReportAnswer{
				Answer: _mapAns,
				Right:  _pa.WasRight,
			})
		}

		players = append(players, &ReportPlayer{
			Player:  _mapPl,
			Answers: ra,
		})
	}

	return &ReportGame{
		Code:    game.Code,
		Players: players,
	}, nil

}
