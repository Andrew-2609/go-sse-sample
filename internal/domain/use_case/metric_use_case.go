package use_case

import (
	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
	"github.com/Andrew-2609/go-sse-sample/internal/domain/enum"
	"github.com/Andrew-2609/go-sse-sample/internal/presentation/dto"
	"github.com/Andrew-2609/go-sse-sample/pkg/sse"
	"github.com/google/uuid"
)

type MetricUseCase struct {
	metricRepository entity.MetricRepository
	sseHub           *sse.SSEHub
}

func NewMetricUseCase(metricRepository entity.MetricRepository) *MetricUseCase {
	return &MetricUseCase{
		metricRepository: metricRepository,
		sseHub:           sse.GetSSEHub(),
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

	response := dto.NewCreateMetricResponseDTO(createdMetric)

	u.sseHub.Broadcast <- sse.NewEvent(enum.EventTypeMetricCreated, response)

	return response, nil
}

func (u *MetricUseCase) GetMetricByID(id uuid.UUID) (dto.GetMetricByIDResponseDTO, error) {
	metric, err := u.metricRepository.GetMetricByID(id)

	if err != nil {
		return dto.GetMetricByIDResponseDTO{}, err
	}

	return dto.NewGetMetricByIDResponseDTO(metric), nil
}
