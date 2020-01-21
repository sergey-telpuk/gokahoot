package repositories

import (
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
		return err
	}

	return nil
}

func (r TestRepository) FindOne(query interface{}, args ...interface{}) (*models.Test, error) {
	var test models.Test

	if err := r.db.GetConn().Where(query, args).First(&test).Error; err != nil {

		return nil, err
	}

	return &test, nil
}

func (r TestRepository) Update(m *models.Test) (*models.Test, error) {
	if err := r.db.GetConn().Save(&m).Error; err != nil {

		return nil, err
	}

	return m, nil
}

func (r TestRepository) Delete(query interface{}, args ...interface{}) error {
	if err := r.db.GetConn().Where(query, args...).Delete(models.Question{}).Error; err != nil {
		return err
	}

	return nil
}

func (r TestRepository) FindAll() ([]*models.Test, error) {
	var tests []*models.Test

	if err := r.db.GetConn().Find(&tests).Limit(1000).Error; err != nil {
		return nil, err
	}

	return tests, nil
}
