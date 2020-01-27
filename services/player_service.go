package services

import (
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
)

const ContainerNamePlayerService = "ContainerNamePlayerService"

type PlayerService struct {
	r *repositories.PlayerRepository
}

func (s *PlayerService) CreateNewPlayer(playerUUID guuid.UUID, gameID int, name string) error {
	model := &models.Player{
		UUID:   playerUUID.String(),
		Name:   name,
		GameID: gameID,
	}

	return s.r.Create(model)
}

func (s *PlayerService) FindByUuid(id string) (*models.Player, error) {
	return s.r.FindOne("players.uuid = ?", id)
}

func (s *PlayerService) FindByID(id int) (*models.Player, error) {
	return s.r.FindOne("players.id = ?", id)
}

func (s *PlayerService) UpdateByUUID(m *models.Player) (*models.Player, error) {
	return s.r.Update(m)
}

func (s *PlayerService) DeleteByUUIDs(id ...string) error {
	return s.r.Delete("players.UUID IN (?)", id)
}

func (s *PlayerService) DeleteByIDs(id ...int) error {
	return s.r.Delete("players.id IN (?)", id)
}

func (s *PlayerService) FindAll() ([]*models.Player, error) {
	return s.r.FindAll()
}

func (s *PlayerService) FindQuestionBelongToGame(id int) ([]*models.Player, error) {
	return s.r.FindQuestionBelongToGame(id)
}

func InitPlayerService(r *repositories.PlayerRepository) *PlayerService {
	return &PlayerService{
		r: r,
	}
}
