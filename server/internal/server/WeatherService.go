package server

import (
	"context"
	"fmt"
	"time"

	"github.com/klados/weather_monitor/internal/model"
	"github.com/klados/weather_monitor/internal/repository"
)

type WeatherService struct {
	repo *repository.WeatherRepository
}

func NewWeatherService(repo *repository.WeatherRepository) *WeatherService {
	return &WeatherService{repo: repo}
}

func (s *WeatherService) SaveSensorWeather(ctx context.Context, sensorName string, data model.SensorWeather) error {
	if data.RecordedAt.IsZero() {
		data.RecordedAt = time.Now().UTC()
	}

	if err := s.repo.SaveSensorWeather(ctx, sensorName, data); err != nil {
		return fmt.Errorf("weather service: %w", err)
	}

	return nil
}
