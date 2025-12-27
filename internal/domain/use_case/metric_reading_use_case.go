package use_case

import (
	"time"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
	"github.com/Andrew-2609/go-sse-sample/internal/domain/enum"
	"github.com/Andrew-2609/go-sse-sample/internal/presentation/dto"
	"github.com/Andrew-2609/go-sse-sample/pkg/sse"
	"github.com/google/uuid"
)

type MetricReadingUseCase struct {
	metricRepository        entity.MetricRepository
	metricReadingRepository entity.MetricReadingRepository
	sseHub                  *sse.SSEHub
}

func NewMetricReadingUseCase(metricRepository entity.MetricRepository, metricReadingRepository entity.MetricReadingRepository, sseHub *sse.SSEHub) *MetricReadingUseCase {
	return &MetricReadingUseCase{
		metricRepository:        metricRepository,
		metricReadingRepository: metricReadingRepository,
		sseHub:                  sseHub,
	}
}

func (u *MetricReadingUseCase) CreateMetricReading(metricReadingDTO dto.CreateMetricReadingRequestDTO) (dto.CreateMetricReadingResponseDTO, error) {
	metricID, err := uuid.Parse(metricReadingDTO.MetricID)
	if err != nil {
		return dto.CreateMetricReadingResponseDTO{}, err
	}

	metric, err := u.metricRepository.GetMetricByID(metricID)

	if err != nil {
		return dto.CreateMetricReadingResponseDTO{}, err
	}

	metricReadingID, err := uuid.NewV7()
	if err != nil {
		return dto.CreateMetricReadingResponseDTO{}, err
	}

	var timestamp time.Time

	if metricReadingDTO.HasTimestamp() {
		timestamp, err = time.Parse(time.RFC3339, *metricReadingDTO.Timestamp)
		if err != nil {
			return dto.CreateMetricReadingResponseDTO{}, err
		}
	}

	metricReadingEntity, err := entity.NewMetricReading(metricReadingID, metric.ID, metricReadingDTO.Value, &timestamp)
	if err != nil {
		return dto.CreateMetricReadingResponseDTO{}, err
	}

	metricReading, err := u.metricReadingRepository.CreateMetricReading(metricReadingEntity)
	if err != nil {
		return dto.CreateMetricReadingResponseDTO{}, err
	}

	response := dto.NewCreateMetricReadingResponseDTO(metricReading)

	u.sseHub.Broadcast <- sse.NewEvent(enum.EventTypeMetricReadingCreated, response)

	return response, nil
}
