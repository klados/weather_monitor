package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/klados/weather_monitor/internal/service"
)

type Weather struct {
	DB             *firestore.Client
	WeatherService service.WeatherService
}

func (we *Weather) WeatherNow(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")

	if location == "" {
		http.Error(w, "Location header is required", http.StatusBadRequest)
		return
	}

	weatherData, err := we.WeatherService.GetWeatherLastHour(r.Context(), location)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weatherData); err != nil {
		http.Error(w, "Failed to encode weather data", http.StatusInternalServerError)
		return
	}
}

func (we *Weather) HistoricalWeather(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	timespanInDays := r.URL.Query().Get("timespanInDays")

	if location == "" || timespanInDays == "" {
		http.Error(w, "Location and timespan query parameters are required", http.StatusBadRequest)
		return
	}

	timespan, err := strconv.ParseUint(timespanInDays, 10, 32)
	if err != nil {
		http.Error(w, "Invalid timespan parameter", http.StatusBadRequest)
		return
	}

	if timespan >= 60 {
		http.Error(w, "Timespan must be less than 60 days", http.StatusBadRequest)
		return
	}

	data, err := we.WeatherService.GetHistoricalWeatherData(r.Context(), location, uint(timespan))
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode weather data", http.StatusInternalServerError)
		return
	}
}
