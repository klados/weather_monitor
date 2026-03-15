package application

import (
	"net/http"

	"firebase.google.com/go/v4/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/klados/weather_monitor/handler"
)

func loadRoutes(fireDb *db.Client) *chi.Mux {
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

func loadApiRoutes(router chi.Router, fireDb *db.Client) {
	weatherHandler := &handler.Weather{DB: fireDb}

	router.Get("/now", weatherHandler.WeatherNow)
}

func loadEmbeddedRoutes(router chi.Router, fireDb *db.Client) {
	sensorHandler := &handler.SensorReceiver{DB: fireDb}

	router.Post("/weather", sensorHandler.SensorData)
}
