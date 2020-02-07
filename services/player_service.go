package services

import (
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
)

const ContainerNamePlayerService = "ContainerNamePlayerService"

type PlayerService struct {
	rpl *repositories.PlayerRepository
	rpa *repositories.PlayerAnswerRepository
}

func (s PlayerService) CreateNewPlayer(playerUUID guuid.UUID, gameID int, name string) error {
	model := &models.Player{
		UUID:   playerUUID.String(),
		Name:   name,
		GameID: gameID,
	}

	return s.rpl.Create(model)
}

func (s PlayerService) CreateNewPlayerAnswer(playerID int, gameID int, questionID int, answerID int, wasRight bool) error {
	model := &models.PlayerAnswer{
		PlayerID:   playerID,
		GameID:     gameID,
		QuestionID: questionID,
		AnswerID:   answerID,
		WasRight:   wasRight,
	}

	return s.rpa.Create(model)
}

func (s PlayerService) FindByUuid(uuid string) (*models.Player, error) {
	return s.rpl.FindOne("players.uuid = ?", uuid)
}

func (s PlayerService) FindPlayerAnswersByGameAndQuestion(game models.Game, question models.Question) ([]*models.PlayerAnswer, error) {
	return s.rpa.Find("player_answers.game_id = ? AND player_answers.question_id = ?", game.ID, question.ID)
}

func (s PlayerService) FindOnePlayerAnswerByGameAndQuestion(game models.Game, player models.Player, question models.Question) (*models.PlayerAnswer, error) {
	return s.rpa.FindOne(
		"player_answers.game_id = ? AND player_answers.question_id = ? AND player_answers.player_id = ?",
		game.ID,
		question.ID,
		player.ID,
	)
}

func (s PlayerService) GetPlayerByUUID(uuid string) (*models.Player, error) {
	m, err := s.FindByUuid(uuid)

	if err != nil {
		return nil, err
	}

	if m == nil {
		return nil, errors.New(fmt.Sprintf("Player model error: %s", "not a such player"))
	}

	return m, nil
}

func (s PlayerService) FindByID(id int) (*models.Player, error) {
	return s.rpl.FindOne("players.id = ?", id)
}

func (s PlayerService) UpdateByUUID(m *models.Player) (*models.Player, error) {
	return s.rpl.Update(m)
}

func (s PlayerService) DeleteByUUIDs(id ...string) error {
	return s.rpl.Delete("players.UUID IN (?)", id)
}

func (s PlayerService) DeleteByIDs(id ...int) error {
	return s.rpl.Delete("players.id IN (?)", id)
}

func (s PlayerService) FindAll() ([]models.Player, error) {
	return s.rpl.FindAll()
}

func (s PlayerService) FindPlayersBelongToGame(id int) ([]models.Player, error) {
	return s.rpl.FindQuestionBelongToGame(id)
}

func InitPlayerService(pl *repositories.PlayerRepository, rpa *repositories.PlayerAnswerRepository) *PlayerService {
	return &PlayerService{
		rpl: pl,
		rpa: rpa,
	}
}
