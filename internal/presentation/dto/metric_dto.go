package dto

import (
	"time"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
)

type CreateMetricRequestDTO struct {
	Name           string `json:"name" binding:"required"`
	InputFrequency string `json:"input_frequency" binding:"required,min=0"`
}

type CreateMetricResponseDTO struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	InputFrequency time.Duration `json:"input_frequency"`
}

func NewCreateMetricResponseDTO(metric entity.Metric) CreateMetricResponseDTO {
	return CreateMetricResponseDTO{
		ID:             metric.ID.String(),
		Name:           metric.Name,
		InputFrequency: metric.InputFrequency,
	}
}

type GetMetricByIDResponseDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewGetMetricByIDResponseDTO(metric entity.Metric) GetMetricByIDResponseDTO {
	return GetMetricByIDResponseDTO{
		ID:   metric.ID.String(),
		Name: metric.Name,
	}
}
