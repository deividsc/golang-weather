package weathertoday

import (
	"context"
	"encoding/json"
	"golang-weather/src/mocks"
	"golang-weather/src/services"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockWeatherService struct {
	response services.WeatherData
	err      error
}

func (s *MockWeatherService) GetData(ctx context.Context, date time.Time) (services.WeatherData, error) {
	return s.response, s.err
}

func TestWeatherHandler(t *testing.T) {
	mockResponse := services.WeatherData{
		Date:        time.Now().Format(time.DateOnly),
		Temperature: 10.0,
		State:       services.Rainy,
	}
	dateOld := time.Now().Add(time.Duration(-1) * time.Hour).Format(time.DateOnly)
	testCases := []struct {
		desc       string
		date       string
		statusCode int
		data       WeatherResponse
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
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/"+tC.date, nil)
			w := httptest.NewRecorder()

			s := MockWeatherService{
				response: mockResponse,
			}
			sut, _ := NewWeatherHandler(&s, mocks.MockLogger)

			sut.GetToday(w, req)

			assert.Equal(t, w.Code, tC.statusCode)

			b := WeatherResponse{}
			err := json.NewDecoder(w.Body).Decode(&b)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tC.data, b)

		})
	}
}
