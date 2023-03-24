package external_services

import (
	"encoding/base64"
	"github.com/Niromash/niromash-api/api"
	"github.com/Niromash/niromash-api/utils/environment"
	"github.com/goccy/go-json"
	"net/http"
	"time"
)

var _ api.WakatimeService = (*wakatimeService)(nil)

type wakatimeService struct {
	service                         *ExternalServicesService
	wakatimeBaseUrl, wakatimeApiKey string
}

func newWakatimeService() *wakatimeService {
	return &wakatimeService{}
}

func (e *wakatimeService) Init(service *ExternalServicesService) error {
	e.service = service
	e.wakatimeBaseUrl = "https://wakatime.com/api/v1"
	e.wakatimeApiKey = environment.GetWakatimeApiKey()

	return nil
}

func (e *wakatimeService) GetTotalDevTime() (resp *api.TotalDevTimeResponse[api.Duration], err error) {
	response, err := e.doWakatimeRequest(e.buildWakatimeRequest("/users/current/all_time_since_today"))
	if err != nil {
		return
	}
	defer response.Body.Close()

	var totalDevTimeResponse api.TotalDevTimeResponse[float64]
	err = json.NewDecoder(response.Body).Decode(&totalDevTimeResponse)
	if err != nil {
		return
	}
	resp = &api.TotalDevTimeResponse[api.Duration]{
		Data: api.TotalDevTimeResponseData[api.Duration]{
			Range:        api.TotalDevTimeResponseDataRange{Start: totalDevTimeResponse.Data.Range.Start},
			TotalSeconds: api.Duration(totalDevTimeResponse.Data.TotalSeconds) * api.Duration(time.Second),
		}}

	return
}

func (e *wakatimeService) GetBestDevTimeDay() (api.Duration, error) {
	response, err := e.doWakatimeRequest(e.buildWakatimeRequest("/users/current/stats/all_time"))
	defer response.Body.Close()

	var totalDevTimeResponse struct {
		Data struct {
			BestDay struct {
				TotalSeconds float64 `json:"total_seconds"`
			} `json:"best_day"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&totalDevTimeResponse)
	if err != nil {
		return 0, err
	}

	return api.Duration(totalDevTimeResponse.Data.BestDay.TotalSeconds) * api.Duration(time.Second), nil
}

func (e *wakatimeService) buildWakatimeRequest(endpoint string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, e.wakatimeBaseUrl+endpoint, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(e.wakatimeApiKey)))
	return request, nil
}

func (e *wakatimeService) doWakatimeRequest(r *http.Request, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (e *wakatimeService) GetLastTodayHeartbeat() (api.Duration, error) {
	response, err := e.doWakatimeRequest(e.buildWakatimeRequest("/users/current/heartbeats?date=" + time.Now().Format("2006-01-02")))
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	var heartbeatsResponse api.HeartbeatsResponse[float64]
	err = json.NewDecoder(response.Body).Decode(&heartbeatsResponse)
	if err != nil {
		return 0, err
	}

	if len(heartbeatsResponse.Data) == 0 {
		return 0, nil
	}

	return api.Duration(heartbeatsResponse.Data[len(heartbeatsResponse.Data)-1].Time), nil
}

func (e *wakatimeService) GetTodayHeartbeats() (resp *api.HeartbeatsResponse[api.Duration], err error) {
	response, err := e.doWakatimeRequest(e.buildWakatimeRequest("/users/current/heartbeats?date=" + time.Now().Format("2006-01-02")))
	if err != nil {
		return
	}
	defer response.Body.Close()

	var heartbeatsResponse api.HeartbeatsResponse[float64]
	err = json.NewDecoder(response.Body).Decode(&heartbeatsResponse)
	if err != nil {
		return
	}
	resp = &api.HeartbeatsResponse[api.Duration]{
		Data: []api.HeartbeatsResponseData[api.Duration]{},
	}
	for _, datum := range heartbeatsResponse.Data {
		resp.Data = append(resp.Data, api.HeartbeatsResponseData[api.Duration]{
			Time: api.Duration(datum.Time),
		})
	}

	return
}
