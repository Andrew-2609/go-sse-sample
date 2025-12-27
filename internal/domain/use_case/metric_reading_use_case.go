package use_case

import (
	"time"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
	"github.com/Andrew-2609/go-sse-sample/internal/presentation/dto"
	"github.com/google/uuid"
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

func (u *MetricReadingUseCase) CreateMetricReading(metricReadingDTO dto.CreateMetricReadingRequestDTO) (dto.CreateMetricReadingResponseDTO, error) {
	metric, err := u.metricRepository.GetMetricByID(metricReadingDTO.MetricID)

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

	return dto.NewCreateMetricReadingResponseDTO(metricReading), nil
}
