package graphql

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/services"
)

func (r *mutationResolver) CreateNewTest(ctx context.Context, input NewTest) (*Test, error) {
	uuid := guuid.New()
	service := r.Di.Container.Get(services.ContainerNameTestService).(*services.TestService)

	if err := service.CreateNewTest(uuid, input.Name, guuid.New()); err != nil {
		return nil, err
	}

	test, err := service.FindByUuid(uuid.String())

	if err != nil {
		return nil, err
	}

	return mapTest(test)
}

func (r *mutationResolver) UpdateTestByUUIDs(ctx context.Context, input []*UpdateTest) ([]*Test, error) {
	var result []*Test
	testService := r.Di.Container.Get(services.ContainerNameTestService).(*services.TestService)
	questionService := r.Di.Container.Get(services.ContainerNameQuestionService).(*services.QuestionService)
	answerService := r.Di.Container.Get(services.ContainerNameAnswerService).(*services.AnswerService)

	for _, iTest := range input {
		mTest, err := testService.FindByUuid(iTest.UUID)

		if err != nil {
			return nil, err
		}

		mTest.Name = iTest.Name

		if _, err := testService.UpdateByUUID(mTest); err != nil {
			return nil, err
		}

		for _, iQuestion := range iTest.Questions {
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

			for _, iAnswer := range iQuestion.Answers {
				mAnswer, err := answerService.FindByID(iAnswer.ID)

				if err != nil {
					return nil, err
				}

				if mAnswer.QuestionID != mQuestion.ID {
					return nil, errors.New(fmt.Sprintf("Question:%s dosent belong to testID: %d", mQuestion.UUID, mAnswer.ID))
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
			}
		}

		mapped, err := mapTest(mTest)
		if err != nil {
			return nil, err
		}

		result = append(result, mapped)
	}

	return result, nil
}

func (r *mutationResolver) DeleteTestByID(ctx context.Context, ids []int) (bool, error) {
	service := r.Di.Container.Get(services.ContainerNameTestService).(*services.TestService)
	if err := service.DeleteByIDs(ids...); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) DeleteTestByUUID(ctx context.Context, ids []string) (bool, error) {
	service := r.Di.Container.Get(services.ContainerNameTestService).(*services.TestService)
	if err := service.DeleteByUUIDs(ids...); err != nil {
		return false, err
	}
	return true, nil
}
