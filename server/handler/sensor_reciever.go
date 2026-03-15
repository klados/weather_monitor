package handler

import (
	"fmt"
	"io"
	"net/http"

	"firebase.google.com/go/v4/db"
)

type SensorReceiver struct {
	DB *db.Client
}

func (sr *SensorReceiver) SensorData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	fmt.Println("received sensor data:", string(body))
}
