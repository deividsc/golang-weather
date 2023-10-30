package todayweather

import (
	"encoding/json"
	"fmt"
	"golang-weather/internal/core/services"
	"log"
	"net/http"
	"time"
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

type WeatherData struct {
	Date           string  `json:"date"`
	MaxTemperature float32 `json:"maxTemperature"`
	MinTemperature float32 `json:"minTemperature"`
	Description    string  `json:"description"`
}

type WeatherResponse struct {
	StatusCode int
	Message    string
	Data       WeatherData `json:"data"`
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

	var response WeatherResponse

	if data, err := h.service.GetData(r.Context(), date); err != nil {
		switch err.(type) {
		case services.EmptyTemperaturesError:
			response = WeatherResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Empty temperature data for selected date",
			}
		default:
			http.Error(w, "Something was wrong!", http.StatusInternalServerError)
		}
	} else {
		response = WeatherResponse{
			StatusCode: http.StatusOK,
			Message:    "Weather data",
			Data: WeatherData{
				Date:           date.Format(time.DateOnly),
				MaxTemperature: data.MaxTemperature,
				MinTemperature: data.MinTemperature,
				Description:    fmt.Sprint(data.Description),
			},
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Fatalf("Error encoding weather: %v", err)
	}
}
