package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Metric struct {
	ID             uuid.UUID
	Name           string
	InputFrequency time.Duration
}

func NewMetric(id uuid.UUID, name string, inputFrequency time.Duration) (Metric, error) {
	metric := Metric{
		ID:             id,
		Name:           name,
		InputFrequency: inputFrequency,
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

	if m.InputFrequency < 0 {
		return errors.New("input frequency must be greater than or equal to 0")
	}

	return nil
}

type MetricRepository interface {
	CreateMetric(metric Metric) (Metric, error)
	GetMetricByID(id uuid.UUID) (Metric, error)
	GetAllMetrics() ([]Metric, error)
}
