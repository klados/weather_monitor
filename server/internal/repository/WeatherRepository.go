package repository

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/klados/weather_monitor/internal/model"
)

const (
	weatherCollectionName = "weather_monitor_project"
	tempHumCollectionName = "temp_hum"
)

type WeatherRepository struct {
	DB *firestore.Client
}

func NewWeatherRepository(dbClient *firestore.Client) *WeatherRepository {
	return &WeatherRepository{DB: dbClient}
}

func (r *WeatherRepository) SaveSensorWeather(ctx context.Context, sensorName string, data model.SensorWeather) error {
	_, _, err := r.DB.
		Collection(weatherCollectionName).
		Doc(sensorName).
		Collection(tempHumCollectionName).
		Add(ctx, data)

	if err != nil {
		return fmt.Errorf("save sensor weather: %w", err)
	}

	return nil
}
