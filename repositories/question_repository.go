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
func (r QuestionRepository) FindOne(query interface{}, args ...interface{}) (*models.Question, error) {
	var question models.Question

	if err := r.db.GetConn().Where(query, args).First(&question).Error; err != nil {
		return nil, err
	}

	return &question, nil
}

func (r QuestionRepository) Find(query interface{}, args ...interface{}) ([]models.Question, error) {
	var questions []models.Question

	if err := r.db.GetConn().Where(query, args).Find(&questions).Limit(10000).Error; err != nil {
		return nil, err
	}

	return questions, nil
}

func (r QuestionRepository) FindQuestionBelongToTest(id int) ([]*models.Question, error) {
	var questions []*models.Question

	if err := r.db.GetConn().Where("test_id = ?", id).Find(&questions).Limit(10000).Error; err != nil {
		return nil, err
	}

	return questions, nil
}

func (r QuestionRepository) FindAll() ([]*models.Question, error) {
	var questions []*models.Question

	if err := r.db.GetConn().Find(&questions).Limit(10000).Error; err != nil {
		return nil, err
	}

	return questions, nil
}
