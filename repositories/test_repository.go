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

func (r TestRepository) Create(model *models.Test) {
	con := r.db.GetConn()

	con.Create(model)
}
func (r TestRepository) FindOne(query interface{}, args ...interface{}) models.Test {
	var test models.Test

	r.db.GetConn().Where(query, args).First(&test)

	return test
}

func (r TestRepository) FindAll() []models.Test {
	var tests []models.Test

	r.db.GetConn().Find(&tests).Limit(1000)

	return tests
}
