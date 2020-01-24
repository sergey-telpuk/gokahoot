package di

import (
	"errors"
	"fmt"
	"github.com/sarulabs/di"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/repositories"
	"github.com/sergey-telpuk/gokahoot/services"
)

type DI struct {
	Container di.Container
}

func (d *DI) Clean() error {
	if err := d.Container.Clean(); err != nil {
		return errorsDI(err)
	}
	return nil
}

func New() (*DI, error) {
	builder, err := di.NewBuilder()

	if err != nil {
		return nil, errorsDI(err)
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
		return nil, errorsDI(err)
	}

	if err := builder.Add(di.Def{
		Name: repositories.ContainerNameTestRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repositories.InitTestRepository(ctn.Get(db.ContainerName).(*db.Db)), nil
		},
	}); err != nil {
		return nil, errorsDI(err)
	}

	if err := builder.Add(di.Def{
		Name: repositories.ContainerNameQuestionRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repositories.InitQuestionRepository(ctn.Get(db.ContainerName).(*db.Db)), nil
		},
	}); err != nil {
		return nil, errorsDI(err)
	}

	if err := builder.Add(di.Def{
		Name: repositories.ContainerNameAnswerRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repositories.InitAnswerRepository(ctn.Get(db.ContainerName).(*db.Db)), nil
		},
	}); err != nil {
		return nil, errorsDI(err)
	}

	if err := builder.Add(di.Def{
		Name: repositories.ContainerNameGameRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repositories.InitGameRepository(ctn.Get(db.ContainerName).(*db.Db)), nil
		},
	}); err != nil {
		return nil, errorsDI(err)
	}

	if err := builder.Add(di.Def{
		Name: repositories.ContainerNamePlayerRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repositories.InitPlayerRepository(ctn.Get(db.ContainerName).(*db.Db)), nil
		},
	}); err != nil {
		return nil, errorsDI(err)
	}

	if err := builder.Add(di.Def{
		Name: services.ContainerNameTestService,
		Build: func(ctn di.Container) (interface{}, error) {
			return services.InitTestService(ctn.Get(repositories.ContainerNameTestRepository).(*repositories.TestRepository)), nil
		},
	}); err != nil {
		return nil, errorsDI(err)
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
		return nil, errorsDI(err)
	}

	if err := builder.Add(di.Def{
		Name: services.ContainerNameAnswerService,
		Build: func(ctn di.Container) (interface{}, error) {
			return services.InitAnswerService(
				ctn.Get(repositories.ContainerNameAnswerRepository).(*repositories.AnswerRepository),
			), nil
		},
	}); err != nil {
		return nil, errorsDI(err)
	}

	if err := builder.Add(di.Def{
		Name: services.ContainerNameGameService,
		Build: func(ctn di.Container) (interface{}, error) {
			return services.InitGameService(
				ctn.Get(repositories.ContainerNameGameRepository).(*repositories.GameRepository),
			), nil
		},
	}); err != nil {
		return nil, errorsDI(err)
	}

	if err := builder.Add(di.Def{
		Name: services.ContainerNamePlayerService,
		Build: func(ctn di.Container) (interface{}, error) {
			return services.InitPlayerService(
				ctn.Get(repositories.ContainerNamePlayerRepository).(*repositories.PlayerRepository),
			), nil
		},
	}); err != nil {
		return nil, errorsDI(err)
	}

	app := builder.Build()

	return &DI{Container: app}, nil
}

func errorsDI(err error) error {
	return errors.New(fmt.Sprintf("Provide container was error, error %v", err))
}
