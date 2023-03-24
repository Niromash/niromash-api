package services

import (
	"context"
	"fmt"
	"github.com/Niromash/niromash-api/api"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron"
	"strconv"
	"time"
)

var _ api.TasksService = (*TaskService)(nil)

type TaskService struct {
	service api.MainService
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (t *TaskService) Init(service api.MainService) error {
	t.service = service
	return nil
}

func (t *TaskService) Start() error {
	if err := t.CheckWakatimeActivityTask(); err != nil {
		fmt.Printf("Error while checking wakatime activity: %v", err)
	}
	c := cron.New()
	if err := c.AddFunc("@every 2m", func() {
		if err := t.CheckWakatimeActivityTask(); err != nil {
			fmt.Println(err)
		}
	}); err != nil {
		return err
	}

	c.Start()

	return nil
}

func (t *TaskService) Close() error {
	return nil
}

func (t *TaskService) Settings() api.ServiceSettings {
	return api.ServiceSettings{
		Priority:         50,
		MustWaitForStart: false,
	}
}

func (t *TaskService) CheckWakatimeActivityTask() error {
	now := time.Now()
	heartbeat, err := t.service.ExternalServices().Wakatime().GetLastTodayHeartbeat()
	if err != nil {
		return err
	}

	currentTime := now.Unix()

	isActive := (currentTime - int64(heartbeat)) <= (5 * 60)
	isRedisActive, err := t.service.Databases().Redis().GetClient().Base().Get(context.Background(), "personal:states:developing").Bool()
	if err != nil && err != redis.Nil {
		return err
	}

	if isActive != isRedisActive {
		t.service.Databases().Redis().GetClient().Base().Set(context.Background(), "personal:states:developing", strconv.FormatBool(isActive), 0)
	}

	return nil
}
