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
			location := r.Header.Get("location")

			// If there's no location header, skip caching and proceed to the handler
			if location == "" {
				next.ServeHTTP(w, r)
				return
			}

			cacheKey := "weather_" + location

			// Check if we have a cached response for this location
			if cachedResponse, found := c.Get(cacheKey); found {
				w.Header().Set("Content-Type", "application/json")
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

			// Call the actual handler
			next.ServeHTTP(recorder, r)

			// Only cache successful 200 OK responses
			if recorder.status == http.StatusOK {
				c.Set(cacheKey, recorder.body.Bytes(), duration)
			}
		})
	}
}
