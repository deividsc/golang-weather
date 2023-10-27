package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeatherApiResponse_GetHighestTemperature(t *testing.T) {
	testCases := []struct {
		desc         string
		temperatures []float32
		want         float32
	}{
		{
			desc: "if temperatures is nil should return 0.0",
			want: 0.0,
		},
		{
			desc:         "Should return highest temperature",
			temperatures: []float32{0.0, 2.0, 0.5, 13.0, 5.0},
			want:         13.0,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			sut := WeatherAPIResponse{
				Hourly: HouyrlyWeatherAPIData{
					Temperature: tC.temperatures,
				},
			}

			assert.Equal(t, tC.want, sut.GetHighestTemperature())
		})
	}
}

func TestWeatherApiResponse_GetState(t *testing.T) {
	testCases := []struct {
		desc           string
		precipitations []float32
		want           WeatherState
	}{
		{
			desc: "if precipitations is nil should return Sunny state",
			want: Sunny,
		},
		{
			desc:           "If precipitations has at least one bigger than 1.0 shoud return Rainy state",
			precipitations: []float32{0.0, 2.0, 0.5, 13.0, 5.0},
			want:           Rainy,
		},
		{
			desc:           "If precipitations has at least one bigger than 0.0 but no one bigger than 1.0 should return Cloudy state",
			precipitations: []float32{0.0, 0.1, 0.5, 0.3, 0.4},
			want:           Cloudy,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			sut := WeatherAPIResponse{
				Hourly: HouyrlyWeatherAPIData{
					Precipitation: tC.precipitations,
				},
			}

			assert.Equal(t, tC.want, sut.GetState())
		})
	}
}
