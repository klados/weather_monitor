package middleware

import (
	"bytes"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

// responseRecorder captures the HTTP response body and status code
type responseRecorder struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	rr.body.Write(b)
	return rr.ResponseWriter.Write(b)
}

// CacheByLocationMiddleware caches HTTP responses based on the 'location' header
func CacheByLocationMiddleware(c *cache.Cache, duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			location := r.URL.Query().Get("location")

			// If there's no location header, skip caching and proceed to the handler
			if location == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Build cache key based on URL path and query parameters
			var cacheKey string
			if r.URL.Path == "/api/now" {
				cacheKey = "weather_now_" + location
			} else if r.URL.Path == "/api/historicData" {
				timespan := r.URL.Query().Get("timespanInDays")
				if timespan == "" {
					// If timespan is missing for historicData, skip caching
					next.ServeHTTP(w, r)
					return
				}
				cacheKey = "weather_historicData_" + location + "_" + timespan
			} else {
				// For any other path, use a generic key
				cacheKey = "weather_" + location
			}

			// Check if we have a cached response for this location
			if cachedResponse, found := c.Get(cacheKey); found {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Cache", "HIT")
				w.WriteHeader(http.StatusOK)
				w.Write(cachedResponse.([]byte))
				return
			}

			// Create the recorder to capture the handler's response
			recorder := &responseRecorder{
				ResponseWriter: w,
				status:         http.StatusOK, // Default to 200 OK
				body:           bytes.NewBuffer(nil),
			}

			// Set cache MISS header before calling the handler
			recorder.Header().Set("X-Cache", "MISS")

			// Call the actual handler
			next.ServeHTTP(recorder, r)

			// Only cache successful 200 OK responses
			if recorder.status == http.StatusOK {
				c.Set(cacheKey, recorder.body.Bytes(), duration)
			}
		})
	}
}
