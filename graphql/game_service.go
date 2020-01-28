package graphql

import (
	"errors"
	"fmt"
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

func (s *GameService) FindByCode(code string) (*models.Game, error) {
	return s.r.FindOne("games.code = ?", code)
}

func (s *GameService) JoinPlayers(m *models.Game) (*models.Game, error) {
	return s.r.FindPlayers(m)
}

func (s *GameService) AddRelationsQuestionsAndPlayers(m *models.Game) (*models.Game, error) {
	return s.r.AddRelationsQuestionsAndPlayers(m)
}

func (s *GameService) IsWaitingForJoining(code string) (bool, error) {
	m, err := s.r.FindOne("games.code = ?", code)

	if m == nil {
		return false, nil
	}

	return models.GameInWaitingPlayers == m.Status, err
}

func (s *GameService) GetGameByCode(code string) (*models.Game, error) {
	game, err := s.FindByCode(code)

	if err != nil {
		return nil, err
	}

	if game == nil {
		return nil, errors.New(fmt.Sprintf("Game model error: %s", "not a such game"))
	}

	return game, nil
}

func (s *GameService) IsPlayingGame(code string) (bool, error) {
	m, err := s.r.FindOne("games.code = ?", code)

	if m == nil {
		return false, nil
	}

	return models.GameInPlaying == m.Status, err
}

func (s *GameService) IsFinishedGame(code string) (bool, error) {
	m, err := s.r.FindOne("games.code = ?", code)

	if m == nil {
		return false, nil
	}

	return models.GameInFinished == m.Status, err
}

func (s *GameService) FindByID(id int) (*models.Game, error) {
	return s.r.FindOne("games.id = ?", id)
}

func (s *GameService) Update(m *models.Game) (*models.Game, error) {
	return s.r.Update(m)
}

func (s *GameService) DeleteByCODEs(id ...string) error {
	return s.r.Delete("games.code IN (?)", id)
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
