package repositories

import (
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
)

type AnswerRepository struct {
	db *db.Db
}

func InitAnswerRepository(db *db.Db) *AnswerRepository {
	return &AnswerRepository{
		db: db,
	}
}

func (r AnswerRepository) Create(model *models.Answer) {
	con := r.db.GetConn()

	con.Create(model)
}

func (r AnswerRepository) FindByQuestionID(id int) []*models.Answer {
	var answers []*models.Answer
	con := r.db.GetConn()

	con.Where("question_id = ?", id).Find(&answers)

	return answers
}
