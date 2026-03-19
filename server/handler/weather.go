package handler

import (
	"net/http"

	"cloud.google.com/go/firestore"
)

type Weather struct {
	DB *firestore.Client
}

func (we *Weather) WeatherNow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("clear weather"))
}
