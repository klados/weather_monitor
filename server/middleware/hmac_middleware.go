package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/klados/weather_monitor/internal/repository"
)

const allowedSkew = 5 * time.Minute

func HmacMiddleware(db *firestore.Client) func(handler http.Handler) http.Handler {
	authRepo := &repository.AuthorizedMicrocontrollers{DB: db}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			deviceID := r.Header.Get("X-Device-ID")
			timestampStr := r.Header.Get("X-Timestamp")
			signature := r.Header.Get("X-Signature")

			if deviceID == "" || timestampStr == "" || signature == "" {
				http.Error(w, "missing authentication headers", http.StatusUnauthorized)
				return
			}

			ts, err := strconv.ParseInt(timestampStr, 10, 64)
			if err != nil {
				http.Error(w, "invalid timestamp", http.StatusUnauthorized)
				return
			}

			requestTime := time.Unix(ts, 0)
			if time.Since(requestTime) > allowedSkew || time.Until(requestTime) > allowedSkew {
				http.Error(w, "request timestamp is outside allowed window", http.StatusUnauthorized)
				return
			}

			microcontroller, err := authRepo.GetAuthorizedMicrocontrollerByDeviceId(deviceID)
			if err != nil {
				http.Error(w, "unauthorized device", http.StatusUnauthorized)
				return
			}

			if !microcontroller.IsActive {
				http.Error(w, "device is inactive", http.StatusUnauthorized)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "failed to read request body", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(body))

			mac := hmac.New(sha256.New, []byte(microcontroller.HMACSecret))

			// Must match the device-side concatenation order exactly.
			mac.Write([]byte(deviceID))
			mac.Write([]byte(timestampStr))
			mac.Write(body)

			expectedSignature := mac.Sum(nil)
			receivedSignature, err := hex.DecodeString(signature)
			if err != nil {
				http.Error(w, "invalid signature format", http.StatusUnauthorized)
				return
			}
			if !hmac.Equal([]byte(receivedSignature), expectedSignature) {
				http.Error(w, "invalid signature", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
