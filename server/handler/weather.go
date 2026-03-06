package handler

import "net/http"

type Weather struct{}

func (we *Weather) WeatherNow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("clear weather"))
}