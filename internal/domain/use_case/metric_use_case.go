package use_case

import (
	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
	"github.com/Andrew-2609/go-sse-sample/internal/presentation/dto"
	"github.com/google/uuid"
)

type MetricUseCase struct {
	metricRepository entity.MetricRepository
}

func NewMetricUseCase(metricRepository entity.MetricRepository) *MetricUseCase {
	return &MetricUseCase{
		metricRepository: metricRepository,
	}
}

func (u *MetricUseCase) CreateMetric(metricDTO dto.CreateMetricRequestDTO) (dto.CreateMetricResponseDTO, error) {
	metricID, err := uuid.NewV7()

	if err != nil {
		return dto.CreateMetricResponseDTO{}, err
	}

	metricEntity, err := entity.NewMetric(metricID, metricDTO.Name)

	if err != nil {
		return dto.CreateMetricResponseDTO{}, err
	}

	createdMetric, err := u.metricRepository.CreateMetric(metricEntity)

	if err != nil {
		return dto.CreateMetricResponseDTO{}, err
	}

	return dto.NewCreateMetricResponseDTO(createdMetric), nil
}

func (u *MetricUseCase) GetMetricByID(id uuid.UUID) (dto.GetMetricByIDResponseDTO, error) {
	metric, err := u.metricRepository.GetMetricByID(id.String())

	if err != nil {
		return dto.GetMetricByIDResponseDTO{}, err
	}

	return dto.NewGetMetricByIDResponseDTO(metric), nil
}
