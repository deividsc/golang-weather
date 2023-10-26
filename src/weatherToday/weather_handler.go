package weathertoday

import (
	"encoding/json"
	"golang-weather/src/services"
	"log"
	"net/http"
	"time"
)

type IWeatherHandler interface {
	GetToday(w http.ResponseWriter, r *http.Request)
}

type WeatherHandler struct {
	service services.IWeatherService
	logger  *log.Logger
}

func NewWeatherHandler(service services.IWeatherService, logger *log.Logger) (*WeatherHandler, error) {
	return &WeatherHandler{
		service: service,
		logger:  logger,
	}, nil
}

type WeatherResponse struct {
	StatusCode int
	Message    string
	Data       services.WeatherData
}

func (h *WeatherHandler) GetToday(w http.ResponseWriter, r *http.Request) {
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
	data, err := h.service.GetData(r.Context(), date.Format(time.DateOnly))
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
