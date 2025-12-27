package repository

import (
	"fmt"
	"sync"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
	"github.com/google/uuid"
)

type MetricInMemoryRepository struct {
	mu      sync.Mutex
	metrics map[uuid.UUID]entity.Metric
}

var _ entity.MetricRepository = (*MetricInMemoryRepository)(nil)

func NewMetricInMemoryRepository() *MetricInMemoryRepository {
	return &MetricInMemoryRepository{
		mu:      sync.Mutex{},
		metrics: make(map[uuid.UUID]entity.Metric),
	}
}

func (r *MetricInMemoryRepository) CreateMetric(metric entity.Metric) (entity.Metric, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics[metric.ID] = metric
	return metric, nil
}

func (r *MetricInMemoryRepository) GetMetricByID(id uuid.UUID) (entity.Metric, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	metric, ok := r.metrics[id]
	if !ok {
		return entity.Metric{}, fmt.Errorf("metric %s not found", id)
	}
	return metric, nil
}
