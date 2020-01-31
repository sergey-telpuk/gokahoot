package services

import (
	"context"
)

func (r *testResolver) Questions(ctx context.Context, obj *Test) ([]*Question, error) {
	var rQuestion []*Question
	service := r.Di.Container.Get(ContainerNameQuestionService).(*QuestionService)

	mQuestion, err := service.FindQuestionBelongToTest(obj.ID)
	if err != nil {
		return nil, err
	}

	for _, q := range mQuestion {
		mapped, _ := mapQuestion(*q)
		rQuestion = append(rQuestion, mapped)
	}

	return rQuestion, nil
}

func (r *queryResolver) Tests(ctx context.Context) ([]*Test, error) {
	var rTests []*Test

	service := r.Di.Container.Get(ContainerNameTestService).(*TestService)

	mTests, err := service.FindAll()

	if err != nil {
		return nil, err
	}

	for _, test := range mTests {
		mapped, _ := mapTest(*test)
		rTests = append(rTests, mapped)
	}

	return rTests, nil
}

func (r *queryResolver) TestByID(ctx context.Context, id int) (*Test, error) {
	service := r.Di.Container.Get(ContainerNameTestService).(*TestService)

	mTest, err := service.FindByID(id)
	if err != nil {
		return nil, err
	}

	return mapTest(*mTest)
}

func (r *queryResolver) TestByUUID(ctx context.Context, id string) (*Test, error) {
	service := r.Di.Container.Get(ContainerNameTestService).(*TestService)

	mTest, err := service.GetTestByUUID(id)

	if err != nil {
		return nil, err
	}

	return mapTest(*mTest)
}
