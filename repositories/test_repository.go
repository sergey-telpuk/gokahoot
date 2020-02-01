package repositories

import (
	"errors"
	"fmt"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
)

const ContainerNameTestRepository = "ContainerNameTestRepository"

type TestRepository struct {
	db *db.Db
}

func InitTestRepository(db *db.Db) *TestRepository {
	return &TestRepository{
		db: db,
	}
}

func (r TestRepository) Create(model *models.Test) error {
	con := r.db.GetConn()

	if err := con.Create(model).Error; err != nil {
		return errorTest(err)
	}

	return nil
}

func (r TestRepository) FindOne(query interface{}, args ...interface{}) (*models.Test, error) {
	var model models.Test

	if err := r.db.GetConn().Where(query, args).First(&model).Error; err != nil {

		return nil, errorTest(err)
	}

	return &model, nil
}

func (r TestRepository) Update(m *models.Test) (*models.Test, error) {
	if err := r.db.GetConn().Save(&m).Error; err != nil {

		return nil, errorTest(err)
	}

	return m, nil
}

func (r TestRepository) Delete(query interface{}, args ...interface{}) error {
	if err := r.db.GetConn().Where(query, args...).Delete(models.Test{}).Error; err != nil {
		return errorTest(err)
	}

	return nil
}

func (r TestRepository) FindAll() ([]*models.Test, error) {
	var _models []*models.Test

	if err := r.db.GetConn().
		Limit(10000).
		Find(&_models).Error; err != nil {
		return nil, errorTest(err)
	}

	return _models, nil
}

func errorTest(err error) error {
	return errors.New(fmt.Sprintf("Test model error: %s", err))
}
