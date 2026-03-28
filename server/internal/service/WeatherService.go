package service

import (
	"context"
	"fmt"
	"math"
	"sort"
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

func (s *WeatherService) GetWeatherLastHour(ctx context.Context, sensorName string) (model.SensorWeather, error) {
	weatherData, err := s.repo.GetWeatherDataSpecificRange(ctx, sensorName, time.Hour)

	if err != nil {
		return model.SensorWeather{}, fmt.Errorf("weather service: %w", err)
	}

	if len(weatherData) == 0 {
		return model.SensorWeather{}, fmt.Errorf("weather service: no data found for sensor %s", sensorName)
	}

	var sumTemp, sumHum float64
	temps := make([]float64, 0, len(weatherData))
	hums := make([]float64, 0, len(weatherData))

	for _, data := range weatherData {
		sumTemp += data.Temperature
		sumHum += data.Humidity
		temps = append(temps, data.Temperature)
		hums = append(hums, data.Humidity)
	}

	n := float64(len(weatherData))
	avgTemp := sumTemp / n
	avgHum := sumHum / n

	sort.Float64s(temps)
	sort.Float64s(hums)

	var medianTemp, medianHum float64
	mid := len(weatherData) / 2
	if len(weatherData)%2 == 0 {
		medianTemp = (temps[mid-1] + temps[mid]) / 2.0
		medianHum = (hums[mid-1] + hums[mid]) / 2.0
	} else {
		medianTemp = temps[mid]
		medianHum = hums[mid]
	}

	// Use median if average is skewed by outliers (e.g. diff > 5.0)
	finalTemp := avgTemp
	if math.Abs(avgTemp-medianTemp) > 5.0 {
		finalTemp = medianTemp
	}
	finalTemp = math.Round(finalTemp*100) / 100

	finalHum := avgHum
	if math.Abs(avgHum-medianHum) > 5.0 {
		finalHum = medianHum
	}
	finalHum = math.Round(finalHum*100) / 100

	return model.SensorWeather{
		Temperature: finalTemp,
		Humidity:    finalHum,
		RecordedAt:  time.Now().UTC(),
	}, nil
}
