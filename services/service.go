package services

import (
	"github.com/Niromash/niromash-api/api"
	"github.com/Niromash/niromash-api/services/database"
	"github.com/Niromash/niromash-api/services/external_services"
	"github.com/Niromash/niromash-api/services/http"
	"sort"
)

var _ api.MainService = (*MainService)(nil)

type MainService struct {
	database         *database.DatabaseService
	http             *http.HttpService
	projects         *ProjectService
	externalServices *external_services.ExternalServicesService
	stats            *StatsService
	task             *TaskService
	messages         *MessageService
	users            *UsersService
}

func NewMainService() *MainService {
	return &MainService{
		database:         database.NewDatabaseService(),
		http:             http.NewHttpService(),
		projects:         NewProjectService(),
		externalServices: external_services.NewExternalServicesService(),
		stats:            NewStatsService(),
		task:             NewTaskService(),
		messages:         NewMessageService(),
		users:            NewUsersService(),
	}
}

func (m *MainService) Init() error {
	for _, initializer := range m.getServicesInitializer() {
		if err := initializer.Init(m); err != nil {
			return err
		}
	}
	return nil
}

func (m *MainService) Start(errCh chan error) {
	go func() {
		for _, service := range m.getServicesStarter() {
			if service.Settings().MustWaitForStart {
				err := service.Start()
				if err != nil {
					errCh <- err
				}
			} else {
				go func(service api.ServiceStarter) {
					err := service.Start()
					if err != nil {
						errCh <- err
					}
				}(service)
			}
		}
	}()
}

func (m *MainService) Close() error {
	for _, starter := range m.getServicesStarter() {
		if err := starter.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (m *MainService) getServicesInitializer() []api.ServiceInitializer {
	return []api.ServiceInitializer{m.database, m.http, m.projects, m.externalServices, m.stats, m.task, m.messages, m.users}
}

func (m *MainService) getServicesStarter() []api.ServiceStarter {
	starters := []api.ServiceStarter{m.database, m.http, m.task}
	sort.SliceStable(starters, func(i, j int) bool {
		// 0 is the highest priority
		return starters[i].Settings().Priority < starters[j].Settings().Priority
	})
	return starters
}

func (m *MainService) Databases() api.DatabaseService {
	return m.database
}

func (m *MainService) Http() api.HttpService {
	return m.http
}

func (m *MainService) Projects() api.ProjectService {
	return m.projects
}

func (m *MainService) ExternalServices() api.ExternalServicesService {
	return m.externalServices
}

func (m *MainService) Stats() api.StatsService {
	return m.stats
}

func (m *MainService) Tasks() api.TasksService {
	return m.task
}

func (m *MainService) Messages() api.MessagesService {
	return m.messages
}

func (m *MainService) Users() api.UsersService {
	return m.users
}
