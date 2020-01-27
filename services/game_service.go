package services

import (
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
)

const ContainerNameGameService = "ContainerNameGameService"

type GameService struct {
	r *repositories.GameRepository
}

func (s *GameService) CreateNewGame(testID int, code guuid.UUID) error {
	model := &models.Game{
		Code:   code.String(),
		TestID: testID,
	}

	return s.r.Create(model)
}

func (s *GameService) FindByUuid(id string) (*models.Game, error) {
	return s.r.FindOne("games.code = ?", id)
}

func (s *GameService) FindByID(id int) (*models.Game, error) {
	return s.r.FindOne("games.id = ?", id)
}

func (s *GameService) UpdateByUUID(m *models.Game) (*models.Game, error) {
	return s.r.Update(m)
}

func (s *GameService) DeleteByUUIDs(id ...string) error {
	return s.r.Delete("games.UUID IN (?)", id)
}

func (s *GameService) DeleteByIDs(id ...int) error {
	return s.r.Delete("games.id IN (?)", id)
}

func (s *GameService) FindAll() ([]*models.Game, error) {
	return s.r.FindAll()
}

func InitGameService(r *repositories.GameRepository) *GameService {
	return &GameService{
		r: r,
	}
}
