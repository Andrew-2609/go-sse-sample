package repository

import (
	"fmt"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
	"github.com/google/uuid"
)

type MetricInMemoryRepository struct {
	metrics map[uuid.UUID]entity.Metric
}

var _ entity.MetricRepository = (*MetricInMemoryRepository)(nil)

func NewMetricInMemoryRepository() *MetricInMemoryRepository {
	return &MetricInMemoryRepository{
		metrics: make(map[uuid.UUID]entity.Metric),
	}
}

func (r *MetricInMemoryRepository) CreateMetric(metric entity.Metric) (entity.Metric, error) {
	r.metrics[metric.ID] = metric
	return metric, nil
}

func (r *MetricInMemoryRepository) GetMetricByID(id uuid.UUID) (entity.Metric, error) {
	metric, ok := r.metrics[id]
	if !ok {
		return entity.Metric{}, fmt.Errorf("metric %s not found", id)
	}
	return metric, nil
}
