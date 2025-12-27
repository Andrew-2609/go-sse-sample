package dto

import "github.com/Andrew-2609/go-sse-sample/internal/domain/entity"

type CreateMetricRequestDTO struct {
	Name string `json:"name" binding:"required"`
}

type CreateMetricResponseDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewCreateMetricResponseDTO(metric entity.Metric) CreateMetricResponseDTO {
	return CreateMetricResponseDTO{
		ID:   metric.ID.String(),
		Name: metric.Name,
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
