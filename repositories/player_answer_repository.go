package repositories

import (
	"errors"
	"fmt"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
)

const ContainerNamePlayerAnswerRepository = "ContainerNamePlayerAnswerRepository"

type PlayerAnswerRepository struct {
	db *db.Db
}

func InitPlayerAnswerRepository(db *db.Db) *PlayerAnswerRepository {
	return &PlayerAnswerRepository{
		db: db,
	}
}

func (r PlayerAnswerRepository) Create(model *models.PlayerAnswer) error {
	con := r.db.GetConn()

	if err := con.Create(model).Error; err != nil {
		return errorPlayerAnswerRepository(err)
	}

	return nil
}

func (r PlayerAnswerRepository) FindOne(query interface{}, args ...interface{}) (*models.PlayerAnswer, error) {
	var model models.PlayerAnswer

	if err := r.db.GetConn().Preload("Player").
		Preload("Game").
		Preload("Question").
		Preload("Answer").
		Where(query, args).First(&model).Error; err != nil {

		return nil, errorPlayerAnswerRepository(err)
	}

	return &model, nil
}

func (r PlayerAnswerRepository) FindPlayers(m *models.PlayerAnswer) (*models.PlayerAnswer, error) {

	if err := r.db.GetConn().Preload("Players").
		Find(&m).Error; err != nil {

		return nil, errorPlayerAnswerRepository(err)
	}

	return m, nil
}
func (r PlayerAnswerRepository) AddRelationsQuestionsAndPlayers(m *models.PlayerAnswer) (*models.PlayerAnswer, error) {

	if err := r.db.GetConn().Preload("Player").
		Preload("Game").
		Preload("Question").
		Preload("Answer").
		Find(&m).Error; err != nil {

		return nil, errorPlayerAnswerRepository(err)
	}

	return m, nil
}

func (r PlayerAnswerRepository) Update(m *models.PlayerAnswer) (*models.PlayerAnswer, error) {
	if err := r.db.GetConn().Save(&m).Error; err != nil {

		return nil, errorPlayerAnswerRepository(err)
	}

	return m, nil
}

func (r PlayerAnswerRepository) Delete(query interface{}, args ...interface{}) error {
	if err := r.db.GetConn().Where(query, args...).Delete(models.PlayerAnswer{}).Error; err != nil {
		return errorPlayerAnswerRepository(err)
	}

	return nil
}

func (r PlayerAnswerRepository) Find(query interface{}, args ...interface{}) ([]*models.PlayerAnswer, error) {
	var _models []*models.PlayerAnswer

	if err := r.db.GetConn().Preload("Player").
		Preload("Game").
		Preload("Question").
		Preload("Answer").
		Where(query, args...).Find(&_models).Limit(10000).Error; err != nil {
		return nil, errorPlayerAnswerRepository(err)
	}

	return _models, nil
}

func (r PlayerAnswerRepository) FindAll() ([]*models.PlayerAnswer, error) {
	var _models []*models.PlayerAnswer

	if err := r.db.GetConn().Preload("Player").
		Preload("Game").
		Preload("Question").
		Preload("Answer").
		Find(&_models).Limit(1000).Error; err != nil {
		return nil, errorPlayerAnswerRepository(err)
	}

	return _models, nil
}

func errorPlayerAnswerRepository(err error) error {
	return errors.New(fmt.Sprintf("PlayerAnswer model error: %s", err))
}
