package services

import (
	"golang-weather/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeatherApiResponse_GetHighestTemperature(t *testing.T) {
	testCases := []struct {
		desc         string
		temperatures []float32
		want         float32
		err          error
	}{
		{
			desc:         "Should return highest temperature",
			temperatures: []float32{0.0, 2.0, 0.5, 13.0, 5.0},
			want:         13.0,
		},
		{
			desc:         "Should return highest temperature of negative values",
			temperatures: []float32{-10.0, -2.0, -0.5, -13.0, -5.0},
			want:         -0.5,
		},
		{
			desc: "Should return an error if temperature slice is empty",
			err:  EmptyTemperaturesError{},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			sut := WeatherAPIResponse{
				Hourly: HouyrlyWeatherAPIData{
					Temperature: tC.temperatures,
				},
			}
			val, err := sut.GetHighestTemperature()

			assert.Equal(t, tC.want, val)
			assert.Equal(t, tC.err, err)
		})
	}
}

func TestWeatherApiResponse_GetLowestTemperature(t *testing.T) {
	testCases := []struct {
		desc         string
		temperatures []float32
		want         float32
		err          error
	}{
		{
			desc:         "Should return lowest temperature",
			temperatures: []float32{1.0, 2.0, 0.5, 13.0, 5.0},
			want:         0.5,
		},
		{
			desc:         "Should return lowest temperature of negative values",
			temperatures: []float32{-10.0, -2.0, -0.5, -13.0, -5.0},
			want:         -13.0,
		},
		{
			desc: "Should return an error if temperature slice is empty",
			err:  EmptyTemperaturesError{},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			sut := WeatherAPIResponse{
				Hourly: HouyrlyWeatherAPIData{
					Temperature: tC.temperatures,
				},
			}
			val, err := sut.GetLowestTemperature()

			assert.Equal(t, tC.want, val)
			assert.Equal(t, tC.err, err)
		})
	}
}

func TestWeatherApiResponse_GetState(t *testing.T) {
	testCases := []struct {
		desc           string
		precipitations []float32
		want           domain.WeatherDescription
	}{
		{
			desc: "if precipitations is nil should return Sunny state",
			want: domain.Sunny,
		},
		{
			desc:           "If precipitations has at least one bigger than 1.0 shoud return Rainy state",
			precipitations: []float32{0.0, 2.0, 0.5, 13.0, 5.0},
			want:           domain.Rainy,
		},
		{
			desc:           "If precipitations has at least one bigger than 0.0 but no one bigger than 1.0 should return Cloudy state",
			precipitations: []float32{0.0, 0.1, 0.5, 0.3, 0.4},
			want:           domain.Cloudy,
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
