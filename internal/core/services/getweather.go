package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type WeatherState string

const (
	Cloudy            WeatherState = "cloudy"
	Sunny             WeatherState = "sunny"
	Rainy             WeatherState = "rainy"
	WEATHER_DATA_CONF              = "temperature,precipitation,cloudcover"
)

type WeatherData struct {
	Date        string
	Temperature float32
	State       WeatherState
}

func (s *WeatherService) GetData(ctx context.Context, date time.Time) (WeatherData, error) {
	response := WeatherData{}
	date_str := date.Format(time.DateOnly)
	url := fmt.Sprintf("%s?latitude=%s&longitude=%s&hourly=%s&start_date=%s&end_date=%s",
		s.apiURL, s.defaultLat, s.defaultLong, WEATHER_DATA_CONF, date_str, date_str)

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
