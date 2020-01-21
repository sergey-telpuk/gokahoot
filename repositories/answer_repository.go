package repositories

import (
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
		return nil, err
	}

	return &m, nil
}

func (r AnswerRepository) Update(m *models.Answer) (*models.Answer, error) {
	if err := r.db.GetConn().Save(&m).Error; err != nil {

		return nil, err
	}

	return m, nil
}

func (r AnswerRepository) Delete(query interface{}, args ...interface{}) error {
	if err := r.db.GetConn().Where(query, args...).Delete(models.Answer{}).Error; err != nil {
		return err
	}

	return nil
}

func (r AnswerRepository) FindByQuestionID(id int) ([]*models.Answer, error) {
	var answers []*models.Answer
	con := r.db.GetConn()

	if err := con.Where("question_id = ?", id).Find(&answers).Error; err != nil {
		return nil, err
	}

	return answers, nil
}
