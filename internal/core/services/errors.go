package services

type EmptyTemperaturesError struct{}

func (e EmptyTemperaturesError) Error() string {
	return "empty temperature slice"
}
