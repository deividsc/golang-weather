package todayweather

import (
	"context"
	"encoding/json"
	"golang-weather/internal/core/domain"
	"golang-weather/internal/core/services"
	"golang-weather/internal/mocks"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockWeatherService struct {
	response domain.Weather
	err      error
}

func (s *MockWeatherService) GetData(ctx context.Context, date time.Time) (domain.Weather, error) {
	return s.response, s.err
}

func TestWeatherHandler(t *testing.T) {
	date := time.Now()
	mockResponse := WeatherData{
		Date:           date.Format(time.DateOnly),
		MaxTemperature: 10.0,
		MinTemperature: 5.0,
		Description:    "rainy",
	}
	dateOld := time.Now().Add(time.Duration(-1) * time.Hour).Format(time.DateOnly)
	testCases := []struct {
		desc       string
		date       string
		statusCode int
		data       WeatherResponse
		errService error
	}{
		{
			desc:       "Empty date should response a weather",
			statusCode: 200,
			data: WeatherResponse{
				StatusCode: 200,
				Message:    "Weather data",
				Data:       mockResponse,
			},
		},
		{
			desc:       "Path with date should be parsed",
			statusCode: 200,
			date:       dateOld,
			data: WeatherResponse{
				StatusCode: 200,
				Message:    "Weather data",
				Data:       mockResponse,
			},
		},
		{
			desc:       "If Service response with error should change status code",
			statusCode: 200,
			data: WeatherResponse{
				StatusCode: 400,
				Message:    "Empty temperature data for selected date",
			},
			errService: services.EmptyTemperaturesError{},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/"+tC.date, nil)
			w := httptest.NewRecorder()

			s := MockWeatherService{
				response: domain.Weather{
					Date:           date,
					MaxTemperature: mockResponse.MaxTemperature,
					MinTemperature: mockResponse.MinTemperature,
					Description:    domain.Rainy,
				},
			}
			s.err = tC.errService
			sut := NewWeatherHandler(&s, mocks.MockLogger)

			sut.GetWeatherData(w, req)

			assert.Equal(t, tC.statusCode, w.Code)

			b := WeatherResponse{}
			err := json.NewDecoder(w.Body).Decode(&b)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tC.data, b)

		})
	}
}
