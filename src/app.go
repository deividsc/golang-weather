package main

import (
	"golang-weather/src/services"
	weathertoday "golang-weather/src/weatherToday"
	"log"
	"net/http"
)

func main() {
	service := services.NewWeatherService(http.DefaultClient, services.API_WEATHER_URL)
	handler := weathertoday.NewWeatherHandler(service, log.Default())
	http.HandleFunc("/", handler.GetWeatherData)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
