package services

import (
	"context"
	"net/http"
	"time"
)

type IWeatherService interface {
	GetData(ctx context.Context, date time.Time) (WeatherData, error)
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

func (s *WeatherAPIResponse) GetState() WeatherState {
	state := Sunny
	for _, v := range s.Hourly.Precipitation {
		if v > 0.0 && state == Sunny {
			state = Cloudy
		}
		if v > 1.0 {
			return Rainy
		}
	}
	return state
}

func (s *WeatherAPIResponse) GetHighestTemperature() float32 {
	var highest float32
	for _, v := range s.Hourly.Temperature {
		if v > highest {
			highest = v
		}
	}
	return highest
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
