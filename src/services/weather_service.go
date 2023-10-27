package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	API_WEATHER_URL   = "https://api.open-meteo.com/v1/forecast"
	BCN_LAT           = "41.3926679"
	BCN_LONG          = "2.1401891"
	WEATHER_DATA_CONF = "temperature,precipitation,cloudcover"
)

type WeatherState string

const (
	Cloudy WeatherState = "cloudy"
	Sunny  WeatherState = "sunny"
	Rainy  WeatherState = "rainy"
)

type WeatherData struct {
	Date        string
	Temperature float32
	State       WeatherState
}

type IWeatherService interface {
	GetData(ctx context.Context, date time.Time) (WeatherData, error)
}

type WeatherService struct {
	client *http.Client
	apiURL string
}

func NewWeatherService(client *http.Client, apiUrl string) *WeatherService {
	return &WeatherService{client: client, apiURL: apiUrl}
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

func (s *WeatherService) GetData(ctx context.Context, date time.Time) (WeatherData, error) {
	response := WeatherData{}
	date_str := date.Format(time.DateOnly)
	url := fmt.Sprintf("%s?latitude=%s&longitude=%s&hourly=%s&start_date=%s&end_date=%s",
		s.apiURL, BCN_LAT, BCN_LONG, WEATHER_DATA_CONF, date_str, date_str)

	r, err := s.client.Get(url)
	if err != nil {
		return response, err
	}
	defer r.Body.Close()

	var apiData []byte
	apiData, err = io.ReadAll(r.Body)
	if err != nil {
		return response, err
	}
	data := WeatherAPIResponse{}
	err = json.Unmarshal(apiData, &data)
	if err != nil {
		return response, err
	}

	return WeatherData{
		Date:        date.Format(time.DateOnly),
		Temperature: data.GetHighestTemperature(),
		State:       data.GetState(),
	}, nil
}
