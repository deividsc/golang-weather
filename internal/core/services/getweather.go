package services

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-weather/internal/core/domain"
	"io"
	"time"
)

const (
	WEATHER_DATA_CONF = "temperature,precipitation,cloudcover"
)

func (s *WeatherService) GetData(ctx context.Context, date time.Time) (domain.Weather, error) {
	response := domain.Weather{}
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
	highT, err := data.GetHighestTemperature()
	if err != nil {
		return response, err
	}
	lowT, err := data.GetLowestTemperature()
	if err != nil {
		return response, err
	}

	return domain.Weather{
		Date:           date,
		MaxTemperature: highT,
		MinTemperature: lowT,
		Description:    data.GetState(),
	}, nil
}
