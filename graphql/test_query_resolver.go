package graphql

import (
	"context"
	"errors"
	"fmt"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/services"
)

func (r *queryResolver) Tests(ctx context.Context) ([]*Test, error) {
	var mTests []models.Test
	var rTests []*Test

	if err := r.Di.Invoke(func(s *services.TestService) error {
		var err error
		mTests, err = s.FindAll()
		if err != nil {
			return err
		}

		return nil

	}); err != nil {
		return nil, errors.New(fmt.Sprintf("Provide container was error, error %v", err))
	}

	for _, test := range mTests {
		rTests = append(rTests, &Test{
			ID:   test.ID,
			Code: test.Code,
			Name: test.Name,
		})
	}

	return rTests, nil
}

func (r *queryResolver) TestByID(ctx context.Context, id int) (*Test, error) {
	var mTest *models.Test

	if err := r.Di.Invoke(func(s *services.TestService) error {
		var err error
		mTest, err = s.FindByID(id)
		if err != nil {
			return err
		}

		return nil

	}); err != nil {
		return nil, errors.New(fmt.Sprintf("Provide container was error, error %v", err))
	}

	return mapTest(mTest), nil
}

func (r *queryResolver) TestByUUID(ctx context.Context, id string) (*Test, error) {
	var mTest *models.Test

	if err := r.Di.Invoke(func(s *services.TestService) error {
		var err error
		mTest, err = s.FindByUuid(id)
		if err != nil {
			return err
		}

		return nil

	}); err != nil {
		return nil, errors.New(fmt.Sprintf("Provide container was error, error %v", err))
	}

	return mapTest(mTest), nil
}

func mapTest(m *models.Test) *Test {

	return &Test{
		ID:   m.ID,
		UUID: m.UUID,
		Code: m.Code,
		Name: m.Name,
	}
}
