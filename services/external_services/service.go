package external_services

import (
	"github.com/Niromash/niromash-api/api"
)

var _ api.ExternalServicesService = (*ExternalServicesService)(nil)

type ExternalServicesService struct {
	service  api.MainService
	wakatime *wakatimeService
	github   *githubService
}

func NewExternalServicesService() *ExternalServicesService {
	return &ExternalServicesService{wakatime: newWakatimeService(), github: newGithubService()}
}

func (e *ExternalServicesService) Init(service api.MainService) error {
	e.service = service
	if err := e.wakatime.Init(e); err != nil {
		return err
	}
	if err := e.github.Init(e); err != nil {
		return err
	}
	return nil
}

func (e *ExternalServicesService) Wakatime() api.WakatimeService {
	return e.wakatime
}

func (e *ExternalServicesService) Github() api.GithubService {
	return e.github
}
