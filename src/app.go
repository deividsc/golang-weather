package main

import (
	"golang-weather/src/services"
	weathertoday "golang-weather/src/weatherToday"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	apiWeatherURL := os.Getenv("API_WEATHER_URL")
	defaultLat := os.Getenv("DEFAULT_LATITUDE")
	defaultLong := os.Getenv("DEFAULT_LONGITUDE")
	service := services.NewWeatherService(http.DefaultClient, apiWeatherURL, defaultLat, defaultLong)
	handler := weathertoday.NewWeatherHandler(service, log.Default())
	http.HandleFunc("/", handler.GetWeatherData)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
