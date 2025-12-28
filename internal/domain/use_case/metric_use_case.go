package use_case

import (
	"time"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
	"github.com/Andrew-2609/go-sse-sample/internal/domain/enum"
	"github.com/Andrew-2609/go-sse-sample/internal/presentation/dto"
	"github.com/Andrew-2609/go-sse-sample/pkg/sse"
	"github.com/google/uuid"
)

type MetricUseCase struct {
	metricRepository        entity.MetricRepository
	metricReadingRepository entity.MetricReadingRepository
	sseHub                  *sse.SSEHub
}

func NewMetricUseCase(metricRepository entity.MetricRepository, metricReadingRepository entity.MetricReadingRepository) *MetricUseCase {
	return &MetricUseCase{
		metricRepository:        metricRepository,
		metricReadingRepository: metricReadingRepository,
		sseHub:                  sse.GetSSEHub(),
	}
}

func (u *MetricUseCase) CreateMetric(metricDTO dto.CreateMetricRequestDTO) (dto.CreateMetricResponseDTO, error) {
	metricID, err := uuid.NewV7()

	if err != nil {
		return dto.CreateMetricResponseDTO{}, err
	}

	inputFrequency, err := time.ParseDuration(metricDTO.InputFrequency)
	if err != nil {
		return dto.CreateMetricResponseDTO{}, err
	}

	metricEntity, err := entity.NewMetric(metricID, metricDTO.Name, inputFrequency)

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

	return dto.NewGetMetricByIDResponseDTO(metric, nil), nil
}

type GetAllMetricsOptions struct {
	WithReadings bool
}

func (u *MetricUseCase) GetAllMetrics(options GetAllMetricsOptions) ([]dto.GetMetricByIDResponseDTO, error) {
	metrics, err := u.metricRepository.GetAllMetrics()
	if err != nil {
		return nil, err
	}

	response := make([]dto.GetMetricByIDResponseDTO, 0, len(metrics))
	for _, metric := range metrics {
		var readings []entity.MetricReading
		var err error

		if options.WithReadings {
			readings, err = u.metricReadingRepository.GetAllReadingsByMetricID(metric.ID)
			if err != nil {
				return nil, err
			}
		}

		response = append(response, dto.NewGetMetricByIDResponseDTO(metric, readings))
	}

	return response, nil
}
