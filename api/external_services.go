package api

import (
	"github.com/goccy/go-json"
	"github.com/google/go-github/v50/github"
	"strconv"
	"sync"
	"time"
)

type ExternalServicesService interface {
	ServiceInitializer
	Wakatime() WakatimeService
	Github() GithubService
}

type WakatimeService interface {
	GetTotalDevTime() (*TotalDevTimeResponse[Duration], error)
	GetBestDevTimeDay() (Duration, error)
	GetTodayHeartbeats() (*HeartbeatsResponse[Duration], error)
	GetLastTodayHeartbeat() (Duration, error)
}
type ListRepositoriesCallback func(repositories []*github.Repository, wg *sync.WaitGroup)

type GithubService interface {
	ListRepositoriesConcurrent(callback ListRepositoriesCallback) error
}

type Duration time.Duration

type HeartbeatsResponse[T float64 | Duration] struct {
	Data []HeartbeatsResponseData[T] `json:"data"`
}

type HeartbeatsResponseData[T float64 | Duration] struct {
	Time T `json:"time"`
}

type TotalDevTimeResponse[T float64 | Duration] struct {
	Data TotalDevTimeResponseData[T] `json:"data"`
}

type TotalDevTimeResponseData[T float64 | Duration] struct {
	Range        TotalDevTimeResponseDataRange `json:"range"`
	TotalSeconds T                             `json:"total_seconds"`
}

type TotalDevTimeResponseDataRange struct {
	Start time.Time `json:"start"`
}

type RepositoriesStored struct {
	PublicRepositories      int `json:"publicRepositories"`
	PrivateRepositories     int `json:"privateRepositories"`
	PublicOwnedRepositories int `json:"publicOwnedRepositories"`
}

func (t TotalDevTimeResponse[T]) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TotalDevTimeResponse[T]) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

func (r RepositoriesStored) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *RepositoriesStored) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}

func (d Duration) MarshalBinary() ([]byte, error) {
	return []byte(strconv.Itoa(int(d))), nil
}

func (d *Duration) UnmarshalBinary(data []byte) error {
	duration, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	*d = Duration(duration)
	return nil
}
