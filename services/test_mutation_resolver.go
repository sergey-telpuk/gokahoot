package services

import (
	"context"
	guuid "github.com/google/uuid"
)

func (r *mutationResolver) CreateNewTest(ctx context.Context, input NewTest) (*Test, error) {
	uuid := guuid.New()
	service := r.Di.Container.Get(ContainerNameTestService).(*TestService)

	if err := service.CreateNewTest(uuid, input.Name, guuid.New()); err != nil {
		return nil, err
	}

	test, err := service.GetTestByUUID(uuid.String())

	if err != nil {
		return nil, err
	}

	return mapTest(test)
}

func (r *mutationResolver) UpdateTestByUUIDs(ctx context.Context, input []*UpdateTest) ([]*Test, error) {
	var result []*Test
	testService := r.Di.Container.Get(ContainerNameTestService).(*TestService)

	for _, iTest := range input {
		mTest, err := testService.GetTestByUUID(iTest.UUID)

		if err != nil {
			return nil, err
		}

		mTest.Name = iTest.Name

		if _, err := testService.UpdateByUUID(mTest); err != nil {
			return nil, err
		}

		if _, err := r.UpdateQuestionsByUUIDs(ctx, mTest.UUID, iTest.Questions); err != nil {
			return nil, err
		}

		mapped, err := mapTest(mTest)
		if err != nil {
			return nil, err
		}

		result = append(result, mapped)
	}

	return result, nil
}

func (r *mutationResolver) DeleteTestByID(ctx context.Context, ids []int) (*Status, error) {
	service := r.Di.Container.Get(ContainerNameTestService).(*TestService)
	if err := service.DeleteByIDs(ids...); err != nil {
		return nil, err
	}
	return &Status{Success: true}, nil
}

func (r *mutationResolver) DeleteTestByUUID(ctx context.Context, ids []string) (*Status, error) {
	service := r.Di.Container.Get(ContainerNameTestService).(*TestService)
	if err := service.DeleteByUUIDs(ids...); err != nil {
		return nil, err
	}
	return &Status{Success: true}, nil
}
