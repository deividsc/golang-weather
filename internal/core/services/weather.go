package services

import (
	"context"
	"golang-weather/internal/core/domain"
	"net/http"
	"time"
)

type IWeatherService interface {
	GetData(ctx context.Context, date time.Time) (domain.Weather, error)
}

type HouyrlyWeatherAPIData struct {
	Time          []string  `json:"time"`
	Temperature   []float32 `json:"temperature"`
	Precipitation []float32 `json:"precipitation"`
}
type WeatherAPIResponse struct {
	Latitude  float32               `json:"latitude"`
	Longitude float32               `json:"longitude"`
	Hourly    HouyrlyWeatherAPIData `json:"hourly"`
}

func (s *WeatherAPIResponse) GetState() domain.WeatherDescription {
	state := domain.Sunny
	for _, v := range s.Hourly.Precipitation {
		if v > 0.0 && state == domain.Sunny {
			state = domain.Cloudy
		}
		if v > 1.0 {
			return domain.Rainy
		}
	}
	return state
}

func (s *WeatherAPIResponse) GetHighestTemperature() (float32, error) {
	tempLen := len(s.Hourly.Temperature)
	if tempLen == 0 {
		return 0, EmptyTemperaturesError{}
	}

	highest := s.Hourly.Temperature[0]
	for _, v := range s.Hourly.Temperature {
		if v > highest {
			highest = v
		}
	}
	return highest, nil
}

func (s *WeatherAPIResponse) GetLowestTemperature() (float32, error) {
	tempLen := len(s.Hourly.Temperature)
	if tempLen == 0 {
		return 0, EmptyTemperaturesError{}
	}
	lowest := s.Hourly.Temperature[0]
	for i := 1; i < tempLen; i++ {
		v := s.Hourly.Temperature[i]
		if v < lowest {
			lowest = v
		}
	}
	return lowest, nil
}

type WeatherService struct {
	client      *http.Client
	apiURL      string
	defaultLat  string
	defaultLong string
}

func NewWeatherService(client *http.Client, apiUrl, lat, long string) *WeatherService {
	return &WeatherService{client: client, apiURL: apiUrl, defaultLat: lat, defaultLong: long}
}
