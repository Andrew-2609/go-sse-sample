package repository

import (
	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
	"github.com/google/uuid"
)

type MetricReadingInMemoryRepository struct {
	metricReadings map[uuid.UUID]entity.MetricReading
}

var _ entity.MetricReadingRepository = (*MetricReadingInMemoryRepository)(nil)

func NewMetricReadingInMemoryRepository() *MetricReadingInMemoryRepository {
	return &MetricReadingInMemoryRepository{
		metricReadings: make(map[uuid.UUID]entity.MetricReading),
	}
}

func (r *MetricReadingInMemoryRepository) CreateMetricReading(metricReading entity.MetricReading) (entity.MetricReading, error) {
	r.metricReadings[metricReading.ID] = metricReading
	return metricReading, nil
}
