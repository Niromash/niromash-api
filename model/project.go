package model

import (
	"github.com/lib/pq"
	"time"
)

type Project struct {
	Id          uint           `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	Link        string         `json:"link"`
	GithubLink  string         `json:"githubLink"`
	Date        time.Time      `json:"date" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Client      string         `json:"client" gorm:"not null;default:'Niromash'"`
	Categories  pq.StringArray `json:"categories" gorm:"type:text[]"`
	TechStack   pq.StringArray `json:"techStack" gorm:"type:text[]"`
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
}

func (p *Project) GetId() uint {
	return p.Id
}

func (p *Project) GetName() string {
	return p.Name
}

func (p *Project) GetDescription() string {
	return p.Description
}

func (p *Project) GetImage() string {
	return p.Image
}

func (p *Project) GetLink() string {
	return p.Link
}

func (p *Project) GetGithubLink() string {
	return p.GithubLink
}

func (p *Project) GetDate() time.Time {
	return p.Date
}

func (p *Project) GetClient() string {
	return p.Client
}

func (p *Project) GetCategories() []string {
	return p.Categories
}

func (p *Project) GetTechStack() []string {
	return p.TechStack
}

func (p *Project) GetImages() []string {
	return p.Images
}
