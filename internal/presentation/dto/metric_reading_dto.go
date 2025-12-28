package dto

import (
	"strings"
	"time"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
)

type CreateMetricReadingRequestDTO struct {
	MetricID  string  `json:"metric_id" binding:"required"`
	Value     float64 `json:"value" binding:"required,min=0"`
	Timestamp *string `json:"timestamp,omitempty"`
}

func (d *CreateMetricReadingRequestDTO) HasTimestamp() bool {
	return d.Timestamp != nil && strings.TrimSpace(*d.Timestamp) != ""
}

type CreateMetricReadingResponseDTO struct {
	ID        string    `json:"id"`
	MetricID  string    `json:"metric_id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func NewCreateMetricReadingResponseDTO(metricReading entity.MetricReading) CreateMetricReadingResponseDTO {
	return CreateMetricReadingResponseDTO{
		ID:        metricReading.ID.String(),
		MetricID:  metricReading.MetricID.String(),
		Value:     metricReading.Value,
		Timestamp: metricReading.Timestamp,
	}
}

type GetMetricReadingResponseDTO struct {
	ID        string    `json:"id"`
	MetricID  string    `json:"metric_id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func NewGetMetricReadingResponseDTO(metricReading entity.MetricReading) GetMetricReadingResponseDTO {
	return GetMetricReadingResponseDTO{
		ID:        metricReading.ID.String(),
		MetricID:  metricReading.MetricID.String(),
		Value:     metricReading.Value,
		Timestamp: metricReading.Timestamp,
	}
}
