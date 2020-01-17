package graphql

//go:generate go run github.com/99designs/gqlgen
import (
	"go.uber.org/dig"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	Di *dig.Container
}

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

type questionResolver struct{ *Resolver }

type testResolver struct{ *Resolver }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
