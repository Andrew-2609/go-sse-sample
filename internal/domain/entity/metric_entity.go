package entity

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Metric struct {
	ID   uuid.UUID
	Name string
}

func NewMetric(id uuid.UUID, name string) (Metric, error) {
	metric := Metric{
		ID:   id,
		Name: name,
	}

	if err := metric.validate(); err != nil {
		return Metric{}, err
	}

	return metric, nil
}

func (m *Metric) validate() error {
	if m.ID == uuid.Nil {
		return errors.New("id is required")
	}

	if m.ID.Version() != 7 {
		return errors.New("id is not a v7 uuid")
	}

	if strings.TrimSpace(m.Name) == "" {
		return errors.New("name is required")
	}

	return nil
}

type MetricRepository interface {
	CreateMetric(metric Metric) (Metric, error)
	GetMetricByID(id uuid.UUID) (Metric, error)
	GetAllMetrics() ([]Metric, error)
}
