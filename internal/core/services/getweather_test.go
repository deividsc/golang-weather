package services

import (
	"context"
	"golang-weather/internal/core/domain"
	"golang-weather/internal/tests"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
	want := domain.Weather{
		Date:           date,
		MaxTemperature: 25.8,
		MinTemperature: 19.1,
		Description:    domain.Rainy,
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

func TestWeatherService_GetData_Error_EmptyTemperature(t *testing.T) {
	apiUrl := "127.0.0.1:9999"
	l, err := net.Listen("tcp", apiUrl)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(tests.ApiWeatherDataEmptyTemp))
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

	_, err = sut.GetData(context.TODO(), date)

	assert.Equal(t, EmptyTemperaturesError{}, err)
}
