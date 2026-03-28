package application

import (
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/klados/weather_monitor/handler"
	"github.com/klados/weather_monitor/internal/repository"
	"github.com/klados/weather_monitor/internal/service"
	appmiddleware "github.com/klados/weather_monitor/middleware"
	"github.com/patrickmn/go-cache"
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
	weatherRepo := repository.NewWeatherRepository(fireStore)
	weatherService := service.NewWeatherService(weatherRepo)

	weatherHandler := &handler.Weather{
		DB:             fireStore,
		WeatherService: *weatherService,
	}

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	memCache := cache.New(5*time.Minute, 10*time.Minute)

	router.With(appmiddleware.CacheByLocationMiddleware(memCache, 5*time.Minute)).Get("/now", weatherHandler.WeatherNow)

}

func loadEmbeddedRoutes(router chi.Router, fireStore *firestore.Client) {
	weatherRepo := repository.NewWeatherRepository(fireStore)
	weatherService := service.NewWeatherService(weatherRepo)
	sensorHandler := &handler.SensorReceiver{Service: weatherService}

	router.With(appmiddleware.HmacMiddleware(fireStore)).Post("/weather", sensorHandler.SensorData)
}
