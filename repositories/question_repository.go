package repositories

import (
	"errors"
	"fmt"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
)

const ContainerNameQuestionRepository = "ContainerNameQuestionRepository"

type QuestionRepository struct {
	db *db.Db
}

func InitQuestionRepository(db *db.Db) *QuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

func (r QuestionRepository) Create(model *models.Question) {
	con := r.db.GetConn()

	con.Create(model)
}

func (r QuestionRepository) FindOne(query interface{}, args ...interface{}) (*models.Question, error) {
	var model models.Question

	if err := r.db.GetConn().Where(query, args).First(&model).Error; err != nil {
		return nil, errorQuestion(err)
	}

	return &model, nil
}

func (r QuestionRepository) Delete(query interface{}, args ...interface{}) error {
	if err := r.db.GetConn().Where(query, args...).Delete(models.Question{}).Error; err != nil {
		return errorQuestion(err)
	}

	return nil
}

func (r QuestionRepository) Find(query interface{}, args ...interface{}) ([]*models.Question, error) {
	var _models []*models.Question

	if err := r.db.GetConn().Where(query, args...).Find(&_models).Limit(10000).Error; err != nil {
		return nil, errorQuestion(err)
	}

	return _models, nil
}

func (r QuestionRepository) FindQuestionBelongToTest(id int) ([]*models.Question, error) {
	return r.Find("test_id = ?", id)
}

func (r QuestionRepository) Update(m *models.Question) (*models.Question, error) {
	if err := r.db.GetConn().Save(&m).Error; err != nil {

		return nil, errorQuestion(err)
	}

	return m, nil
}

func (r QuestionRepository) FindAll() ([]*models.Question, error) {
	var _models []*models.Question

	if err := r.db.GetConn().Find(&_models).Limit(10000).Error; err != nil {
		return nil, errorQuestion(err)
	}

	return _models, nil
}

func errorQuestion(err error) error {
	return errors.New(fmt.Sprintf("Question model error: %s", err))
}
