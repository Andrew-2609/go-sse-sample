package repository

import (
	"fmt"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
)

type MetricInMemoryRepository struct {
	metrics map[string]entity.Metric
}

var _ entity.MetricRepository = (*MetricInMemoryRepository)(nil)

func NewMetricInMemoryRepository() *MetricInMemoryRepository {
	return &MetricInMemoryRepository{
		metrics: make(map[string]entity.Metric),
	}
}

func (r *MetricInMemoryRepository) CreateMetric(metric entity.Metric) (entity.Metric, error) {
	r.metrics[metric.ID.String()] = metric
	return metric, nil
}

func (r *MetricInMemoryRepository) GetMetricByID(id string) (entity.Metric, error) {
	metric, ok := r.metrics[id]
	if !ok {
		return entity.Metric{}, fmt.Errorf("metric %s not found", id)
	}
	return metric, nil
}
