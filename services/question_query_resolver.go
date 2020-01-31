package services

import (
	"context"
)

func (r *questionResolver) Answers(ctx context.Context, obj *Question) ([]*Answer, error) {
	var rAnswer []*Answer
	service := r.Di.Container.Get(ContainerNameQuestionService).(*QuestionService)

	mAnswer, err := service.FindAnswersBelongToQuestion(obj.ID)
	if err != nil {
		return nil, err
	}

	for _, answer := range mAnswer {
		mapped, _ := mapAnswer(*answer)
		rAnswer = append(rAnswer, mapped)
	}

	return rAnswer, nil
}

func (r *queryResolver) QuestionByID(ctx context.Context, id int) (*Question, error) {
	service := r.Di.Container.Get(ContainerNameQuestionService).(*QuestionService)

	mQuestion, err := service.FindByID(id)

	if err != nil {
		return nil, err
	}

	return mapQuestion(*mQuestion)
}

func (r *queryResolver) QuestionByUUID(ctx context.Context, id string) (*Question, error) {
	service := r.Di.Container.Get(ContainerNameQuestionService).(*QuestionService)

	mQuestion, err := service.FindByUuid(id)

	if err != nil {
		return nil, err
	}

	return mapQuestion(*mQuestion)
}
