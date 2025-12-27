package use_case

import (
	"fmt"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
)

type MetricReadingUseCase struct {
	metricRepository        entity.MetricRepository
	metricReadingRepository entity.MetricReadingRepository
}

func NewMetricReadingUseCase(metricRepository entity.MetricRepository, metricReadingRepository entity.MetricReadingRepository) *MetricReadingUseCase {
	return &MetricReadingUseCase{
		metricRepository:        metricRepository,
		metricReadingRepository: metricReadingRepository,
	}
}

func (u *MetricReadingUseCase) CreateMetricReading(metricReading entity.MetricReading) error {
	exists, err := u.metricRepository.MetricExistsById(metricReading.MetricID)

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("metric %s does not exist", metricReading.MetricID)
	}

	return u.metricReadingRepository.CreateMetricReading(metricReading)
}
