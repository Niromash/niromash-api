package api

import (
	"github.com/Niromash/niromash-api/model"
	"gorm.io/gorm"
)

type MainService interface {
	Init() error
	Start(errCh chan error)
	Close() error
	Databases() DatabaseService
	Http() HttpService
	Projects() ProjectService
	ExternalServices() ExternalServicesService
	Stats() StatsService
	Tasks() TasksService
	Messages() MessagesService
	Users() UsersService
}

type ServiceSettings struct {
	MustWaitForStart bool
	Priority         int
}

type ServiceInitializer interface {
	Init(service MainService) error
}

type ServiceStarter interface {
	Start() error
	Close() error
	Settings() ServiceSettings
}

type DatabaseService interface {
	ServiceStarter
	AutoReconnect(database Database[any]) error
	Postgres() Database[*gorm.DB]
	Redis() Database[RedisClient]
}

type Database[T any] interface {
	WaitForStart()
	GetClient() T
	GetName() string
	Connect() error
	Disconnect() error
	Ping() bool
}

type ProjectService interface {
	ServiceInitializer
	GetProject(id uint) (*model.Project, error)
	ListProjects() ([]*model.Project, error)
}

type StatsService interface {
	ServiceInitializer
	GetTotalDevTime() (*TotalDevTimeResponse[Duration], error)
	GetBestDevTimeDay() (Duration, error)
	IsDeveloping() (bool, error)
	GetVisitorCount() (int, error)
	ListRepositories() (*RepositoriesStored, error)
}

type TasksService interface {
	ServiceInitializer
	ServiceStarter
	CheckWakatimeActivityTask() error
}
