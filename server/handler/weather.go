package handler

import (
	"net/http"

	"firebase.google.com/go/v4/db"
)

type Weather struct {
	DB *db.Client
}

func (we *Weather) WeatherNow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("clear weather"))
}
