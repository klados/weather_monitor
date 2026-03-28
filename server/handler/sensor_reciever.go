package handler

import (
	"encoding/json"
	"net/http"

	"github.com/klados/weather_monitor/internal/model"
	"github.com/klados/weather_monitor/internal/service"
)

type SensorReceiver struct {
	Service *service.WeatherService
}

func (sr *SensorReceiver) SensorData(w http.ResponseWriter, r *http.Request) {
	sensorName := r.Header.Get("X-Device-Id")
	if sensorName == "" {
		http.Error(w, "missing X-Device-Id header", http.StatusBadRequest)
		return
	}

	var payload model.SensorWeather

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := sr.Service.SaveSensorWeather(r.Context(), sensorName, payload); err != nil {
		http.Error(w, "failed to save weather data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("weather data saved"))
}
