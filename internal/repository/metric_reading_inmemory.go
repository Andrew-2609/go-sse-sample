package repository

import (
	"sort"
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

func (r *MetricReadingInMemoryRepository) GetLastMetricReadingValue(metricID uuid.UUID) (float64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	metricReadings := make([]entity.MetricReading, 0)

	for _, metricReading := range r.metricReadings {
		if metricReading.MetricID == metricID {
			metricReadings = append(metricReadings, metricReading)
		}
	}

	if len(metricReadings) == 0 {
		return 0, nil
	}

	sort.Slice(metricReadings, func(i, j int) bool {
		return metricReadings[i].Timestamp.After(metricReadings[j].Timestamp)
	})

	return metricReadings[0].Value, nil
}
