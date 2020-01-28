package repositories

import (
	"errors"
	"fmt"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
)

const ContainerNamePlayerRepository = "ContainerNamePlayerRepository"

type PlayerRepository struct {
	db *db.Db
}

func InitPlayerRepository(db *db.Db) *PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}

func (r PlayerRepository) Create(model *models.Player) error {
	con := r.db.GetConn()

	if err := con.Create(model).Error; err != nil {
		return errorPlayer(err)
	}

	return nil
}

func (r PlayerRepository) FindOne(query interface{}, args ...interface{}) (*models.Player, error) {
	var model models.Player

	if err := r.db.GetConn().Preload("Game").
		Preload("Game.Test").
		Where(query, args).First(&model).Error; err != nil {

		return nil, errorPlayer(err)
	}

	return &model, nil
}

func (r PlayerRepository) Update(m *models.Player) (*models.Player, error) {
	if err := r.db.GetConn().Save(&m).Error; err != nil {

		return nil, errorPlayer(err)
	}

	return m, nil
}

func (r PlayerRepository) Delete(query interface{}, args ...interface{}) error {
	if err := r.db.GetConn().Where(query, args...).Delete(models.Player{}).Error; err != nil {
		return errorPlayer(err)
	}

	return nil
}

func (r PlayerRepository) FindAll() ([]*models.Player, error) {
	var _models []*models.Player

	if err := r.db.GetConn().Preload("Game").
		Preload("Game.Test").
		Find(&_models).Limit(1000).Error; err != nil {
		return nil, errorPlayer(err)
	}

	return _models, nil
}

func (r PlayerRepository) Find(query interface{}, args ...interface{}) ([]*models.Player, error) {
	var _models []*models.Player

	if err := r.db.GetConn().Preload("Game").
		Preload("Game.Test").
		Where(query, args...).Find(&_models).Limit(10000).Error; err != nil {
		return nil, errorPlayer(err)
	}

	return _models, nil
}

func (r PlayerRepository) FindQuestionBelongToGame(id int) ([]*models.Player, error) {
	return r.Find("players.game_id = ?", id)
}

func errorPlayer(err error) error {
	return errors.New(fmt.Sprintf("Player model error: %s", err))
}
