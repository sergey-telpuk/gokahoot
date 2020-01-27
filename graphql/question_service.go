package graphql

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
)

const ContainerNameQuestionService = "ContainerNameQuestionService"

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

func (s *QuestionService) FindByUuid(id string) (*models.Question, error) {
	return s.rq.FindOne("uuid = ?", id)
}

func (s *QuestionService) FindByID(id int) (*models.Question, error) {
	return s.rq.FindOne("id = ?", id)
}
func (s *QuestionService) DeleteByIDs(id ...int) error {
	return s.rq.Delete("id IN (?)", id)
}
func (s *QuestionService) DeleteByUUIDs(id ...string) error {
	return s.rq.Delete("uuid IN (?)", id)
}

func (s *QuestionService) FindAll() ([]*models.Question, error) {
	return s.rq.FindAll()
}

func (s *QuestionService) UpdateByUUID(m *models.Question) (*models.Question, error) {
	return s.rq.Update(m)
}

func (s *QuestionService) FindAnswersBelongToQuestion(id int) ([]*models.Answer, error) {
	return s.ra.FindByQuestionID(id)
}

func (s *QuestionService) FindQuestionBelongToTest(id int) ([]*models.Question, error) {
	return s.rq.FindQuestionBelongToTest(id)
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
