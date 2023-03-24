package services

import (
	"context"
	"github.com/Niromash/niromash-api/api"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-github/v50/github"
	"strings"
	"sync"
	"time"
)

var _ api.StatsService = (*StatsService)(nil)

type StatsService struct {
	service api.MainService
}

func NewStatsService() *StatsService {
	return &StatsService{}
}

func (s *StatsService) Init(service api.MainService) error {
	s.service = service
	return nil
}

func (s *StatsService) GetTotalDevTime() (*api.TotalDevTimeResponse[api.Duration], error) {
	var result api.TotalDevTimeResponse[api.Duration]
	if err := s.service.Databases().Redis().GetClient().Base().Get(context.Background(), "personal:stats:totaldevtime").Scan(&result); err != nil {
		if err != redis.Nil {
			return nil, err
		}
		totalDevTime, err := s.service.ExternalServices().Wakatime().GetTotalDevTime()
		if err != nil {
			return nil, err
		}
		if err = s.service.Databases().Redis().GetClient().Base().Set(context.Background(), "personal:stats:totaldevtime", totalDevTime, 5*time.Minute).Err(); err != nil {
			return nil, err
		}
		result = *totalDevTime
	}
	return &result, nil
}

func (s *StatsService) GetBestDevTimeDay() (api.Duration, error) {
	bestDevTimeDayResult, err := s.service.Databases().Redis().GetClient().Base().Get(context.Background(), "personal:stats:bestdevtimeday").Int()
	if err != nil {
		if err != redis.Nil {
			return 0, err
		}
		bestDevTimeDay, err := s.service.ExternalServices().Wakatime().GetBestDevTimeDay()
		if err != nil {
			return 0, err
		}
		if err = s.service.Databases().Redis().GetClient().Base().Set(context.Background(), "personal:stats:bestdevtimeday", bestDevTimeDay, 1*time.Hour).Err(); err != nil {
			return 0, err
		}
		bestDevTimeDayResult = int(bestDevTimeDay)
	}
	return api.Duration(bestDevTimeDayResult), nil
}

func (s *StatsService) IsDeveloping() (bool, error) {
	return s.service.Databases().Redis().GetClient().Base().Get(context.Background(), "personal:states:developing").Bool()
}

func (s *StatsService) GetVisitorCount() (int, error) {
	return s.service.Databases().Redis().GetClient().Base().Get(context.Background(), "personal:states:visitors").Int()
}

func (s *StatsService) ListRepositories() (*api.RepositoriesStored, error) {
	var storedRepositories api.RepositoriesStored
	if err := s.service.Databases().Redis().GetClient().Base().Get(context.Background(), "personal:stats:repositories").Scan(&storedRepositories); err != nil {
		if err != redis.Nil {
			return nil, err
		}
		if err = s.service.ExternalServices().Github().ListRepositoriesConcurrent(func(repositories []*github.Repository, wg *sync.WaitGroup) {
			defer wg.Done()
			for _, repo := range repositories {
				if repo.GetPrivate() {
					storedRepositories.PrivateRepositories++
				} else {
					if strings.EqualFold(repo.GetOwner().GetLogin(), "niromash") {
						storedRepositories.PublicOwnedRepositories++
					}
					storedRepositories.PublicRepositories++
				}
			}
		}); err != nil {
			return nil, err
		}

		err = s.service.Databases().Redis().GetClient().Base().Set(context.Background(), "personal:stats:repositories", storedRepositories, 1*time.Hour).Err()
		if err != nil {
			return nil, err
		}
	}

	return &storedRepositories, nil
}
