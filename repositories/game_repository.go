package repositories

import (
	"errors"
	"fmt"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
)

const ContainerNameGameRepository = "ContainerNameGameRepository"

type GameRepository struct {
	db *db.Db
}

func InitGameRepository(db *db.Db) *GameRepository {
	return &GameRepository{
		db: db,
	}
}

func (r GameRepository) Create(model *models.Game) error {
	con := r.db.GetConn()

	if err := con.Create(model).Error; err != nil {
		return errorGame(err)
	}

	return nil
}

func (r GameRepository) FindOne(query interface{}, args ...interface{}) (*models.Game, error) {
	var model models.Game

	if err := r.db.GetConn().Preload("Test").
		Where(query, args).First(&model).Error; err != nil {

		return nil, errorGame(err)
	}

	return &model, nil
}

func (r *GameRepository) FindPlayers(m *models.Game) (*models.Game, error) {

	if err := r.db.GetConn().Preload("Players").
		Find(&m).Error; err != nil {

		return nil, errorGame(err)
	}

	return m, nil
}
func (r *GameRepository) AddRelationsQuestionsAndPlayers(m *models.Game) (*models.Game, error) {

	if err := r.db.GetConn().Preload("Players").
		Preload("Test").
		Preload("Test.Questions").
		Preload("Test.Questions.Answers").
		Find(&m).Error; err != nil {

		return nil, errorGame(err)
	}

	return m, nil
}

func (r GameRepository) Update(m *models.Game) (*models.Game, error) {
	if err := r.db.GetConn().Save(&m).Error; err != nil {

		return nil, errorGame(err)
	}

	return m, nil
}

func (r GameRepository) Delete(query interface{}, args ...interface{}) error {
	if err := r.db.GetConn().Where(query, args...).Delete(models.Game{}).Error; err != nil {
		return errorGame(err)
	}

	return nil
}

func (r GameRepository) Find(query interface{}, args ...interface{}) ([]*models.Game, error) {
	var _models []*models.Game

	if err := r.db.GetConn().Preload("Players").
		Preload("Test").
		Preload("Test.Questions").
		Preload("Test.Questions.Answers").
		Where(query, args...).Find(&_models).Limit(10000).Error; err != nil {
		return nil, errorPlayer(err)
	}

	return _models, nil
}

func (r GameRepository) FindAll() ([]*models.Game, error) {
	var _models []*models.Game

	if err := r.db.GetConn().Preload("Players").
		Preload("Test").
		Preload("Test.Questions").
		Preload("Test.Questions.Answers").
		Find(&_models).Limit(1000).Error; err != nil {
		return nil, errorGame(err)
	}

	return _models, nil
}

func errorGame(err error) error {
	return errors.New(fmt.Sprintf("Game model error: %s", err))
}
