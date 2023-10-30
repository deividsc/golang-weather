package domain

import "time"

type WeatherDescription string

const (
	Sunny  WeatherDescription = "sunny"
	Rainy  WeatherDescription = "rainy"
	Cloudy WeatherDescription = "cloudy"
)

type Weather struct {
	Date           time.Time
	MaxTemperature float32
	MinTemperature float32
	Description    WeatherDescription
}
