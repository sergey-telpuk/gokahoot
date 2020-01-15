package repositories

import (
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
)

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
func (r QuestionRepository) FindOne(query interface{}, args ...interface{}) models.Question {
	var test models.Question

	r.db.GetConn().Where(query, args).First(&test)

	return test
}

func (r QuestionRepository) FindAll() []models.Question {
	var questions []models.Question

	r.db.GetConn().Find(&questions).Limit(1000)

	return questions
}
