package di

import (
	"github.com/sarulabs/di"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/repositories"
	"github.com/sergey-telpuk/gokahoot/services"
	"log"
)

type DI struct {
	Container di.Container
}

func (d *DI) Clean() {
	if err := d.Container.Clean(); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}
}

func New() *DI {
	builder, err := di.NewBuilder()

	if err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	if err := builder.Add(di.Def{
		Name: db.ContainerName,
		Build: func(ctn di.Container) (interface{}, error) {
			return db.Init()
		},
		Close: func(obj interface{}) error {
			return obj.(*db.Db).GetConn().Close()
		},
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	if err := builder.Add(di.Def{
		Name: repositories.ContainerNameTestRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repositories.InitTestRepository(ctn.Get(db.ContainerName).(*db.Db)), nil
		},
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	if err := builder.Add(di.Def{
		Name: repositories.ContainerNameQuestionRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repositories.InitQuestionRepository(ctn.Get(db.ContainerName).(*db.Db)), nil
		},
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	if err := builder.Add(di.Def{
		Name: repositories.ContainerNameAnswerRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repositories.InitAnswerRepository(ctn.Get(db.ContainerName).(*db.Db)), nil
		},
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	if err := builder.Add(di.Def{
		Name: services.ContainerNameTestService,
		Build: func(ctn di.Container) (interface{}, error) {
			return services.InitTestService(ctn.Get(repositories.ContainerNameTestRepository).(*repositories.TestRepository)), nil
		},
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	if err := builder.Add(di.Def{
		Name: services.ContainerNameQuestionService,
		Build: func(ctn di.Container) (interface{}, error) {
			return services.InitQuestionService(
				ctn.Get(repositories.ContainerNameQuestionRepository).(*repositories.QuestionRepository),
				ctn.Get(repositories.ContainerNameAnswerRepository).(*repositories.AnswerRepository),
				ctn.Get(db.ContainerName).(*db.Db),
			), nil
		},
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	if err := builder.Add(di.Def{
		Name: services.ContainerNameAnswerService,
		Build: func(ctn di.Container) (interface{}, error) {
			return services.InitAnswerService(
				ctn.Get(repositories.ContainerNameAnswerRepository).(*repositories.AnswerRepository),
			), nil
		},
	}); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	app := builder.Build()

	return &DI{Container: app}
}
