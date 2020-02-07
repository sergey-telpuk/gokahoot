package services

import (
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
)

const ContainerNameGameService = "ContainerNameGameService"

type GameService struct {
	rg  *repositories.GameRepository
	rch *repositories.ChatMessageRepository
}

func (s *GameService) CreateNewGame(testID int, code guuid.UUID) error {
	model := &models.Game{
		Code:   code.String(),
		TestID: testID,
	}

	return s.rg.Create(model)
}

func (s *GameService) CreateNewMessageOfChat(uuid guuid.UUID, gameID int, playerID int, message string) error {
	model := &models.ChatMessage{
		UUID:     uuid.String(),
		Message:  message,
		GameID:   gameID,
		PlayerID: playerID,
	}

	return s.rch.Create(model)
}

func (s *GameService) FindByCode(code string) (*models.Game, error) {
	return s.rg.FindOne("games.code = ?", code)
}

func (s *GameService) JoinPlayers(m *models.Game) (*models.Game, error) {
	return s.rg.FindPlayers(m)
}

func (s *GameService) AddRelationsQuestionsAndPlayers(m *models.Game) (*models.Game, error) {
	return s.rg.AddRelationsQuestionsAndPlayers(m)
}

func (s *GameService) IsWaitingForJoining(code string) (bool, error) {
	m, err := s.rg.FindOne("games.code = ?", code)

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

func (s *GameService) GetChatMessageByUUID(uuid string) (*models.ChatMessage, error) {
	message, err := s.rch.FindOne("chat_messages.uuid = ?", uuid)

	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, errors.New(fmt.Sprintf("ChatMessage model error: %s", "not a such message"))
	}

	return message, nil
}

func (s *GameService) FindChatMessagesByGameCode(code string, offset int, limit int, order ChatTimeOrder) ([]*models.ChatMessage, error) {
	messages, err := s.rch.Find(offset, limit, string(order), "games.code = ?", code)

	if err != nil {
		return nil, err
	}

	if messages == nil {
		return nil, errors.New(fmt.Sprintf("ChatMessage model error: %s", "not a such message"))
	}

	return messages, nil
}

func (s *GameService) IsPlayingGame(code string) (bool, error) {
	m, err := s.rg.FindOne("games.code = ?", code)

	if m == nil {
		return false, nil
	}

	return models.GameInPlaying == m.Status, err
}

func (s *GameService) IsFinishedGame(code string) (bool, error) {
	m, err := s.rg.FindOne("games.code = ?", code)

	if m == nil {
		return false, nil
	}

	return models.GameInFinished == m.Status, err
}

func (s *GameService) FindByID(id int) (*models.Game, error) {
	return s.rg.FindOne("games.id = ?", id)
}

func (s *GameService) Update(m *models.Game) (*models.Game, error) {
	return s.rg.Update(m)
}

func (s *GameService) DeleteByCODEs(id ...string) error {
	return s.rg.Delete("games.code IN (?)", id)
}

func (s *GameService) DeleteByIDs(id ...int) error {
	return s.rg.Delete("games.id IN (?)", id)
}

func (s *GameService) FindAll() ([]*models.Game, error) {
	return s.rg.FindAll()
}

func (s *GameService) FindAllWhichAreWaitingForJoining() ([]*models.Game, error) {
	return s.rg.Find("games.status = ?", models.GameInWaitingPlayers)
}

func (s *GameService) FindAllWhichArePlaying() ([]*models.Game, error) {
	return s.rg.Find("games.status = ?", models.GameInPlaying)
}

func InitGameService(rg *repositories.GameRepository, rch *repositories.ChatMessageRepository) *GameService {
	return &GameService{
		rg:  rg,
		rch: rch,
	}
}
