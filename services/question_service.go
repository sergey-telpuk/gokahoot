package services

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
)

type QuestionService struct {
	rq *repositories.QuestionRepository
	ra *repositories.AnswerRepository
	db *db.Db
}

func (s *QuestionService) CreateNewQuestion(
	uuid uuid.UUID,
	testID int,
	text string,
	imgURL *string,
	rightAnswer int,
	answers []*models.Answer,
) error {

	return s.db.GetConn().Transaction(func(tx *gorm.DB) error {
		question := &models.Question{
			UUID:        uuid.String(),
			TestID:      testID,
			Text:        text,
			ImgURL:      imgURL,
			RightAnswer: rightAnswer,
		}

		tx.Create(question).Scan(question)

		if err := tx.Error; err != nil {
			// return any error will rollback
			return err
		}

		for _, answer := range answers {
			answer.QuestionID = question.ID

			if err := tx.Create(answer).Error; err != nil {
				return err
			}
		}

		// return nil will commit
		return nil
	})
}

func (s *QuestionService) FindByUuid(id uuid.UUID) models.Question {
	return s.rq.FindOne("uuid = ?", id.String())
}

func (s *QuestionService) FindAll() []models.Question {
	return s.rq.FindAll()
}

func (s *QuestionService) FindAnswersBelongToQuestion(id int) []*models.Answer {
	return s.ra.FindByQuestionID(id)
}

func InitQuestionService(
	rq *repositories.QuestionRepository,
	ra *repositories.AnswerRepository,
	db *db.Db,
) *QuestionService {
	return &QuestionService{
		rq: rq,
		ra: ra,
		db: db,
	}
}
