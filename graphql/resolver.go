package graphql

//go:generate go run github.com/99designs/gqlgen
import (
	"errors"
	"fmt"
	"github.com/sergey-telpuk/gokahoot/di"
	"github.com/sergey-telpuk/gokahoot/models"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	Di *di.DI
}

type questionResolver struct{ *Resolver }

type testResolver struct{ *Resolver }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

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

func (r *Resolver) Test() TestResolver {
	return &testResolver{r}
}

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

func mapQuestion(m *models.Question) (*Question, error) {

	if m.ID == 0 {
		return nil, errors.New(fmt.Sprintf("Not a such item"))
	}

	return &Question{
		ID:          m.ID,
		UUID:        m.UUID,
		TestID:      m.TestID,
		Text:        m.Text,
		ImgURL:      m.ImgURL,
		RightAnswer: m.RightAnswer,
	}, nil
}

func mapAnswer(m *models.Answer) (*Answer, error) {

	if m.ID == 0 {
		return nil, errors.New(fmt.Sprintf("Not a such item"))
	}

	return &Answer{
		ID:         m.ID,
		Text:       m.Text,
		Sequential: m.Sequential,
		ImgURL:     m.ImgURL,
	}, nil
}

func mapTest(m *models.Test) (*Test, error) {

	if m.ID == 0 {
		return nil, errors.New(fmt.Sprintf("Not a such item"))
	}

	return &Test{
		ID:   m.ID,
		UUID: m.UUID,
		Name: m.Name,
	}, nil

}

func mapGame(m *models.Game) (*Game, error) {

	if m.ID == 0 {
		return nil, errors.New(fmt.Sprintf("Not a such item"))
	}

	return &Game{
		Code:     m.Code,
		TestUUID: m.TestID,
	}, nil

}
