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
		return nil, errors.New(fmt.Sprintf("Creating question was failed, error %v", err))
	}

	question, err := service.FindByUuid(uuid.String())

	if err != nil {
		return nil, err
	}

	return mapQuestion(question)
}

func (r *mutationResolver) DeleteQuestionByID(ctx context.Context, ids []int) (*Status, error) {
	service := r.Di.Container.Get(services.ContainerNameQuestionService).(*services.QuestionService)
	if err := service.DeleteByIDs(ids...); err != nil {
		return nil, err
	}
	return &Status{Success: true}, nil
}

func (r *mutationResolver) DeleteQuestionByUUID(ctx context.Context, ids []string) (*Status, error) {
	service := r.Di.Container.Get(services.ContainerNameQuestionService).(*services.QuestionService)
	if err := service.DeleteByUUIDs(ids...); err != nil {
		return nil, err
	}
	return &Status{Success: true}, nil
}

func (r *mutationResolver) UpdateQuestionsByUUIDs(ctx context.Context, testUUID string, questions []*UpdateQuestion) ([]*Question, error) {
	var result []*Question
	questionService := r.Di.Container.Get(services.ContainerNameQuestionService).(*services.QuestionService)
	testService := r.Di.Container.Get(services.ContainerNameTestService).(*services.TestService)

	mTest, err := testService.FindByUuid(testUUID)

	if err != nil {
		return nil, err
	}

	for _, iQuestion := range questions {
		mQuestion, err := questionService.FindByUuid(iQuestion.UUID)
		if err != nil {
			return nil, err
		}

		if mQuestion.TestID != mTest.ID {
			return nil, errors.New(fmt.Sprintf("Question:%s dosent belong to testID: %s", mQuestion.UUID, mTest.UUID))
		}

		if iQuestion.RightAnswer != nil {
			mQuestion.RightAnswer = *iQuestion.RightAnswer
		}
		if iQuestion.Text != nil {
			mQuestion.Text = *iQuestion.Text
		}
		mQuestion.ImgURL = iQuestion.ImgURL

		if _, err := questionService.UpdateByUUID(mQuestion); err != nil {
			return nil, err
		}
		mapped, err := mapQuestion(mQuestion)

		if err != nil {
			return nil, err
		}

		result = append(result, mapped)

		if _, err := r.UpdateAnswersByIDs(ctx, mQuestion.UUID, iQuestion.Answers); err != nil {
			return nil, err
		}

	}

	return result, nil
}

func (r *mutationResolver) UpdateAnswersByIDs(ctx context.Context, questionUUID string, answers []*UpdateAnswer) ([]*Answer, error) {
	var result []*Answer
	questionService := r.Di.Container.Get(services.ContainerNameQuestionService).(*services.QuestionService)
	answerService := r.Di.Container.Get(services.ContainerNameAnswerService).(*services.AnswerService)

	mQuestion, err := questionService.FindByUuid(questionUUID)

	if err != nil {
		return nil, err
	}

	for _, iAnswer := range answers {
		mAnswer, err := answerService.FindByID(iAnswer.ID)

		if err != nil {
			return nil, err
		}

		if mAnswer.QuestionID != mQuestion.ID {
			return nil, errors.New(fmt.Sprintf("Answer:%s dosent belong to testID: %s", mQuestion.UUID, mQuestion.UUID))
		}

		if iAnswer.Text != nil {
			mAnswer.Text = *iAnswer.Text
		}
		if iAnswer.Sequential != nil {
			mAnswer.Sequential = *iAnswer.Sequential
		}

		mAnswer.ImgURL = iAnswer.ImgURL

		if _, err := answerService.UpdateByUUID(mAnswer); err != nil {
			return nil, err
		}

		mapped, err := mapAnswer(mAnswer)

		result = append(result, mapped)
	}

	return result, nil
}
