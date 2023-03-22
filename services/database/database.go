package database

import (
	"fmt"
	"gorm.io/gorm"
	"niromash-api/api"
	"time"
)

var _ api.DatabaseService = (*DatabaseService)(nil)

type DatabaseService struct {
	postgres *PostgresDatabase
	redis    *RedisDatabase
}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

func (d *DatabaseService) Init(service api.MainService) error {
	d.postgres = &PostgresDatabase{service: service}
	d.redis = &RedisDatabase{service: service}

	return nil
}

func (d *DatabaseService) Start() error {
	if err := d.postgres.Connect(); err != nil {
		return err
	}
	fmt.Printf("Connected to %s\n", d.postgres.GetName())

	if err := d.redis.Connect(); err != nil {
		return err
	}
	fmt.Printf("Connected to %s\n", d.redis.GetName())

	return nil
}

func (d *DatabaseService) Close() error {
	if err := d.postgres.Disconnect(); err != nil {
		return err
	}

	if err := d.redis.Disconnect(); err != nil {
		return err
	}

	return nil
}

func (d *DatabaseService) AutoReconnect(database api.Database[any]) error {
	for {
		if database.Ping() {
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf("Reconnecting to %s...\n", database.GetName())
		if err := database.Connect(); err != nil {
			fmt.Printf("Failed to reconnect to %s: %v", database.GetName(), err)
			time.Sleep(5 * time.Second)
			continue
		}
	}
}

func (d *DatabaseService) Settings() api.ServiceSettings {
	return api.ServiceSettings{
		MustWaitForStart: true,
		Priority:         0,
	}
}

func (d *DatabaseService) Postgres() api.Database[*gorm.DB] {
	return d.postgres
}

func (d *DatabaseService) Redis() api.Database[api.RedisClient] {
	return d.redis
}
