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

// corsMiddleware handles CORS headers and preflight requests
func corsMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set the allowed origin from config
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Device-Id")
			w.Header().Set("Access-Control-Expose-Headers", "X-Cache")

			// Handle preflight OPTIONS requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func loadRoutes(fireDb *firestore.Client, allowedOrigin string) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/api", func(r chi.Router) {
		r.Use(corsMiddleware(allowedOrigin))

		// Catch-all route for OPTIONS requests to ensure the CORS middleware is triggered
		r.Options("/*", func(w http.ResponseWriter, r *http.Request) {})

		loadApiRoutes(r, fireDb)
	})

	router.Route("/embedded", func(r chi.Router) {
		loadEmbeddedRoutes(r, fireDb)
	})

	// Serve static files in production
	loadStaticRoutes(router)

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

	router.With(appmiddleware.CacheByLocationMiddleware(memCache, 5*time.Minute)).Get("/historicalData", weatherHandler.HistoricalWeather)
}

func loadEmbeddedRoutes(router chi.Router, fireStore *firestore.Client) {
	weatherRepo := repository.NewWeatherRepository(fireStore)
	weatherService := service.NewWeatherService(weatherRepo)
	sensorHandler := &handler.SensorReceiver{Service: weatherService}

	router.With(appmiddleware.HmacMiddleware(fireStore)).Post("/weather", sensorHandler.SensorData)
}

func loadStaticRoutes(router *chi.Mux) {
	// Serve static assets from the React build
	staticDir := "./frontend/dist"

	// Serve static files (CSS, JS, images, etc.)
	fileServer := http.FileServer(http.Dir(staticDir))

	router.Handle("/assets/*", fileServer)
	router.Handle("/favicon.ico", fileServer)
	router.Handle("/robots.txt", fileServer)

	// Serve index.html for all other routes (SPA support)
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/index.html")
	})
}
