package repositories

import (
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
)

type TestRepository struct {
	db *db.Db
}

//TODO
func InitTestRepository(db *db.Db) *TestRepository {
	return &TestRepository{
		db: db,
	}
}

func (r TestRepository) Create(model *models.TestModel) {
	con := r.db.GetConn()

	con.Create(model)
}
func (r TestRepository) FindOne(query interface{}, args ...interface{}) models.TestModel {
	var test models.TestModel

	r.db.GetConn().Where(query, args).First(&test)

	return test
}

func (r TestRepository) FindAll() []models.TestModel {
	var tests []models.TestModel

	r.db.GetConn().Find(&tests).Limit(1000)

	return tests
}
