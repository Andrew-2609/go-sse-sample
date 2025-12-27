package repository

import (
	"sync"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
	"github.com/google/uuid"
)

type MetricReadingInMemoryRepository struct {
	mu             sync.Mutex
	metricReadings map[uuid.UUID]entity.MetricReading
}

var _ entity.MetricReadingRepository = (*MetricReadingInMemoryRepository)(nil)

func NewMetricReadingInMemoryRepository() *MetricReadingInMemoryRepository {
	return &MetricReadingInMemoryRepository{
		mu:             sync.Mutex{},
		metricReadings: make(map[uuid.UUID]entity.MetricReading),
	}
}

func (r *MetricReadingInMemoryRepository) CreateMetricReading(metricReading entity.MetricReading) (entity.MetricReading, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metricReadings[metricReading.ID] = metricReading
	return metricReading, nil
}
