package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/klados/weather_monitor/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/api", loadApiRoutes)
	
	router.Route("/embedded", loadEmbeddedRoutes)
	
	return router
}

func loadApiRoutes(router chi.Router) {
	weatherHandler := &handler.Weather{}

	router.Get("/now", weatherHandler.WeatherNow)
}

func loadEmbeddedRoutes(router chi.Router) {
	
}
