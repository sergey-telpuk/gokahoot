package repositories

import (
	"errors"
	"fmt"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
)

const ContainerNameAnswerRepository = "ContainerNameAnswerRepository"

type AnswerRepository struct {
	db *db.Db
}

func InitAnswerRepository(db *db.Db) *AnswerRepository {
	return &AnswerRepository{
		db: db,
	}
}

func (r AnswerRepository) Create(m *models.Answer) error {
	con := r.db.GetConn()

	if err := con.Create(m).Error; err != nil {
		return err
	}

	return nil
}

func (r AnswerRepository) FindOne(query interface{}, args ...interface{}) (*models.Answer, error) {
	var m models.Answer

	if err := r.db.GetConn().Where(query, args).First(&m).Error; err != nil {
		return nil, errorAnswer(err)
	}

	return &m, nil
}

func (r AnswerRepository) Update(m *models.Answer) (*models.Answer, error) {
	if err := r.db.GetConn().Save(&m).Error; err != nil {

		return nil, errorAnswer(err)
	}

	return m, nil
}

func (r AnswerRepository) Delete(query interface{}, args ...interface{}) error {
	if err := r.db.GetConn().Where(query, args...).Delete(models.Answer{}).Error; err != nil {
		return errorAnswer(err)
	}

	return nil
}

func (r AnswerRepository) FindByQuestionID(id int) ([]*models.Answer, error) {
	var _models []*models.Answer
	con := r.db.GetConn()

	if err := con.Where("question_id = ?", id).Find(&_models).Error; err != nil {
		return nil, errorAnswer(err)
	}

	return _models, nil
}

func errorAnswer(err error) error {
	return errors.New(fmt.Sprintf("Answer model error: %s", err))
}
