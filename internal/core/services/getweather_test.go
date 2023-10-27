package services

import (
	"context"
	"golang-weather/internal/tests"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestWeatherService_GetData(t *testing.T) {
	apiUrl := "127.0.0.1:9999"
	l, err := net.Listen("tcp", apiUrl)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(tests.ApiWeatherData))
		if err != nil {
			t.Fatal(err)
		}
	}))
	srv.Listener.Close()
	srv.Listener = l
	srv.Start()
	defer srv.Close()

	sut := NewWeatherService(http.DefaultClient, "http://"+apiUrl, "test", "test")

	date := time.Now()
	want := WeatherData{
		Date:        date.Format(time.DateOnly),
		Temperature: 25.8,
		State:       Rainy,
	}

	data, err := sut.GetData(context.TODO(), date)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, data)
}

func TestWeatherService_GetData_Error(t *testing.T) {
	sut := NewWeatherService(http.DefaultClient, "http://notexists.test", "test", "test")

	date := time.Now()

	_, err := sut.GetData(context.TODO(), date)

	assert.NotNil(t, err)
}
