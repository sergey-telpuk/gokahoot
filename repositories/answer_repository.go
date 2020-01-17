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

func (r AnswerRepository) Create(model *models.Answer) error {
	con := r.db.GetConn()

	if err := con.Create(model).Error; err != nil {
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
