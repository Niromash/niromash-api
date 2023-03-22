package database

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"niromash-api/api"
	"niromash-api/model"
	"niromash-api/utils/environment"
	"time"
)

var _ api.Database[*gorm.DB] = (*PostgresDatabase)(nil)

type PostgresDatabase struct {
	client  *gorm.DB
	service api.MainService
}

func (p *PostgresDatabase) GetClient() *gorm.DB {
	return p.client
}

func (p *PostgresDatabase) GetName() string {
	return "Postgres"
}

func (p *PostgresDatabase) Connect() (err error) {
	p.client, err = gorm.Open(postgres.Open(environment.GetPostgresDSN()), &gorm.Config{})
	if err != nil {
		return
	}
	if err = p.client.AutoMigrate(&model.Project{}, &model.User{}, &model.MessageArgument{}, &model.MessageTranslation{},
		&model.Scope{}, &model.Message{}); err != nil {
		return
	}
	return
}

func (p *PostgresDatabase) Disconnect() error {
	sqlDB, err := p.client.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *PostgresDatabase) Ping() bool {
	db, err := p.client.DB()
	if err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.PingContext(ctx) == nil
}

func (p *PostgresDatabase) WaitForStart() {
	for {
		if p.client != nil && p.Ping() {
			return
		}

		fmt.Printf("Waiting for %s to start...\n", p.GetName())
		time.Sleep(1 * time.Second)
	}
}
