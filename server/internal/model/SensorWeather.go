package model

import "time"

type SensorWeather struct {
	Temperature float64   `json:"temperature" firestore:"temperature"`
	Humidity    float64   `json:"humidity" firestore:"humidity"`
	RecordedAt  time.Time `json:"recorded_at" firestore:"recorded_at"`
}
