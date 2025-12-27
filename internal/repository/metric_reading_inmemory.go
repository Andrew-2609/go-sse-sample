package repository

import (
	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
)

type MetricReadingInMemoryRepository struct {
	metricReadings map[string]entity.MetricReading
}

var _ entity.MetricReadingRepository = (*MetricReadingInMemoryRepository)(nil)

func NewMetricReadingInMemoryRepository() *MetricReadingInMemoryRepository {
	return &MetricReadingInMemoryRepository{
		metricReadings: make(map[string]entity.MetricReading),
	}
}

func (r *MetricReadingInMemoryRepository) CreateMetricReading(metricReading entity.MetricReading) error {
	r.metricReadings[metricReading.ID] = metricReading
	return nil
}
