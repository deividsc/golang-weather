package todayweather

import (
	"encoding/json"
	"golang-weather/internal/core/services"
	"log"
	"net/http"
	"time"
)

const (
	WEATHER_DATA_CONF = "temperature,precipitation,cloudcover"
)

type WeatherHandler struct {
	service services.IWeatherService
	logger  *log.Logger
}

func NewWeatherHandler(service services.IWeatherService, logger *log.Logger) *WeatherHandler {
	return &WeatherHandler{
		service: service,
		logger:  logger,
	}
}

type WeatherResponse struct {
	StatusCode int
	Message    string
	Data       services.WeatherData
}

func (h *WeatherHandler) GetWeatherData(w http.ResponseWriter, r *http.Request) {
	date := time.Now()
	d := r.URL.Path[1:]
	if d != "" {
		var err error
		date, err = time.Parse(time.DateOnly, d)
		if err != nil {
			http.Error(w, "Bad format date! Use yyyy-mm-dd", http.StatusBadRequest)
			return
		}
	}
	data, err := h.service.GetData(r.Context(), date)
	if err != nil {
		http.Error(w, "Something was wrong!", http.StatusInternalServerError)
	}

	response := WeatherResponse{
		StatusCode: http.StatusOK,
		Message:    "Weather data",
		Data:       data,
	}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Fatalf("Error encoding weather: %v", err)
	}
}
