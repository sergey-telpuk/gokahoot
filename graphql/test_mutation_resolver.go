package graphql

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/services"
)

func (r *mutationResolver) CreateNewTest(ctx context.Context, input NewTest) (*Test, error) {
	var test *models.Test
	uuid := guuid.New()

	if err := r.Di.Invoke(func(s *services.TestService) error {
		if err := s.CreateNewTest(uuid, input.Name, guuid.New()); err != nil {
			return err
		}
		var err error
		test, err = s.FindByUuid(uuid.String())

		if err != nil {
			return err
		}

		return nil

	}); err != nil {
		return nil, errors.New(fmt.Sprintf("Provide container was error, error %v", err))
	}

	if test.ID != 0 {
		return &Test{
			ID:        test.ID,
			UUID:      test.UUID,
			Code:      test.Code,
			Name:      test.Name,
			Questions: nil,
		}, nil
	}

	return nil, nil
}
