package services

import (
	"context"
	"net/http"
)

type WeatherState string

const (
	Cloudy WeatherState = "cloudy"
	Sunny  WeatherState = "sunny"
	Rainy  WeatherState = "rainy"
)

type WeatherData struct {
	Date        string
	Temperature string
	State       WeatherState
}

type IWeatherService interface {
	GetData(ctx context.Context, date string) (WeatherData, error)
}

type WeatherService struct {
	client *http.Client
}

func NewWeatherService(client *http.Client) *WeatherService {
	return &WeatherService{client: client}
}

func (s *WeatherService) GetData(ctx context.Context, date string) (WeatherData, error) {
	return WeatherData{}, nil
}
