package repositories

import (
	"errors"
	"fmt"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
)

const ContainerNameChatMessageRepository = "ContainerNameChatMessageRepository"

type ChatMessageRepository struct {
	db *db.Db
}

func InitChatMessageRepository(db *db.Db) *ChatMessageRepository {
	return &ChatMessageRepository{
		db: db,
	}
}

func (r ChatMessageRepository) Create(model *models.ChatMessage) error {
	con := r.db.GetConn()

	if err := con.Create(model).Error; err != nil {
		return errorChatMessageGame(err)
	}

	return nil
}

func (r ChatMessageRepository) FindOne(query interface{}, args ...interface{}) (*models.ChatMessage, error) {
	var model models.ChatMessage

	if err := r.db.GetConn().Preload("Game").
		Preload("Player").
		Preload("Player.Game").
		Where(query, args...).First(&model).Error; err != nil {

		return nil, errorGame(err)
	}

	return &model, nil
}

func (r ChatMessageRepository) Find(offset int, limit int, order string, query interface{}, args ...interface{}) ([]*models.ChatMessage, error) {
	var _models []*models.ChatMessage

	if err := r.db.GetConn().Preload("Game").
		Joins("left join games on games.id = chat_messages.game_id").
		Preload("Player").
		Preload("Player.Game").
		Where(query, args...).
		Order("chat_messages.created_at " + order).
		Limit(limit).
		Offset(offset).
		Find(&_models).Error; err != nil {
		return nil, errorChatMessageGame(err)
	}

	return _models, nil
}

func (r ChatMessageRepository) FindAll(offset int, limit int) ([]*models.ChatMessage, error) {
	var _models []*models.ChatMessage

	if err := r.db.GetConn().Preload("Game").
		Preload("Player").
		Preload("Player.Game").
		Order("chat_messages.created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&_models).Error; err != nil {
		return nil, errorChatMessageGame(err)
	}

	return _models, nil
}

func errorChatMessageGame(err error) error {
	return errors.New(fmt.Sprintf("ChatMessage model error: %s", err))
}
