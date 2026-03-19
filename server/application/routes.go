package application

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/klados/weather_monitor/handler"
	"github.com/klados/weather_monitor/internal/repository"
	"github.com/klados/weather_monitor/internal/server"
)

func loadRoutes(fireDb *firestore.Client) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/api", func(r chi.Router) {
		loadApiRoutes(r, fireDb)
	})

	router.Route("/embedded", func(r chi.Router) {
		loadEmbeddedRoutes(r, fireDb)
	})

	return router
}

func loadApiRoutes(router chi.Router, fireStore *firestore.Client) {
	weatherHandler := &handler.Weather{DB: fireStore}

	router.Get("/now", weatherHandler.WeatherNow)
}

func loadEmbeddedRoutes(router chi.Router, fireStore *firestore.Client) {
	weatherRepo := repository.NewWeatherRepository(fireStore)
	weatherService := server.NewWeatherService(weatherRepo)
	sensorHandler := &handler.SensorReceiver{Service: weatherService}

	router.Post("/weather", sensorHandler.SensorData)
}
