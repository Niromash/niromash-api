package services

import (
	"gorm.io/gorm"
	"niromash-api/api"
	"niromash-api/model"
)

var _ api.ProjectService = (*ProjectService)(nil)

type ProjectService struct {
	service api.MainService
}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

func (p *ProjectService) Init(service api.MainService) error {
	p.service = service
	return nil
}

func (p *ProjectService) GetProject(id uint) (project *model.Project, err error) {
	if err = p.service.Databases().Postgres().GetClient().First(&project, model.Project{Id: id}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = api.ErrProjectNotFound
		}
		return
	}
	return
}

func (p *ProjectService) ListProjects() (projects []*model.Project, err error) {
	if err = p.service.Databases().Postgres().GetClient().Order("id").Find(&projects).Error; err != nil {
		return
	}
	return
}
