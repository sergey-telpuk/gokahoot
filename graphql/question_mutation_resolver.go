package graphql

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/services"
)

func (r *mutationResolver) CreateNewQuestion(ctx context.Context, input NewQuestion) (*Question, error) {
	uuid := guuid.New()
	service := r.Di.Container.Get(services.ContainerNameQuestionService).(*services.QuestionService)

	var answers []*models.Answer
	for _, answer := range input.Answers {
		answers = append(answers, &models.Answer{
			Text:       answer.Text,
			ImgURL:     answer.ImgURL,
			Sequential: answer.Sequential,
		})
	}
	if err := service.CreateNewQuestion(
		uuid,
		input.TestID,
		input.Text,
		input.ImgURL,
		input.RightAnswer,
		answers,
	); err != nil {
		return nil, errors.New(fmt.Sprintf("Craeting question was failed, error %v", err))
	}

	question, err := service.FindByUuid(uuid.String())

	if err != nil {
		return nil, err
	}

	return mapQuestion(question)
}

func (r *mutationResolver) DeleteQuestionByID(ctx context.Context, ids []int) (bool, error) {
	service := r.Di.Container.Get(services.ContainerNameQuestionService).(*services.QuestionService)
	if err := service.DeleteByIDs(ids...); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) DeleteQuestionByUUID(ctx context.Context, ids []string) (bool, error) {
	service := r.Di.Container.Get(services.ContainerNameQuestionService).(*services.QuestionService)
	if err := service.DeleteByUUIDs(ids...); err != nil {
		return false, err
	}
	return true, nil
}
