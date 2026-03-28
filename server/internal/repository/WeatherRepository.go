package repository

import (
	"context"
	"fmt"
	"time"

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

func (r *WeatherRepository) GetWeatherDataSpecificRange(ctx context.Context, sensorName string, timeRange time.Duration) ([]model.SensorWeather, error) {
	startTime := time.Now().UTC().Add(-timeRange)

	iter := r.DB.Collection(weatherCollectionName).Doc(sensorName).Collection(tempHumCollectionName).
		Where("recorded_at", ">=", startTime).
		OrderBy("recorded_at", firestore.Desc).
		Documents(ctx)

	docs, err := iter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %w", err)
	}

	var weatherData []model.SensorWeather
	for _, doc := range docs {
		var data model.SensorWeather
		if err := doc.DataTo(&data); err != nil {
			return nil, fmt.Errorf("failed to decode weather data: %w", err)
		}
		weatherData = append(weatherData, data)
	}

	return weatherData, nil
}
