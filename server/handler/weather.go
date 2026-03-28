package handler

import (
	"encoding/json"
	"net/http"

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
