package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type MetricReading struct {
	ID        uuid.UUID
	MetricID  uuid.UUID
	Value     float64
	Timestamp time.Time
}

func NewMetricReading(id uuid.UUID, metricID uuid.UUID, value float64, timestamp *time.Time) (MetricReading, error) {
	reading := MetricReading{
		ID:       id,
		MetricID: metricID,
		Value:    value,
	}

	if timestamp != nil && !timestamp.IsZero() {
		reading.Timestamp = timestamp.UTC()
	} else {
		now := time.Now().UTC()
		reading.Timestamp = now
	}

	if err := reading.validate(); err != nil {
		return MetricReading{}, err
	}

	return reading, nil
}

func (m *MetricReading) validate() error {
	if m.ID == uuid.Nil {
		return errors.New("id is required")
	}

	if m.ID.Version() != 7 {
		return errors.New("id is not a v7 uuid")
	}

	if m.MetricID == uuid.Nil {
		return errors.New("metric id is required")
	}

	if m.MetricID.Version() != 7 {
		return errors.New("metric id is not a v7 uuid")
	}

	if m.Value <= 0 {
		return errors.New("value must be greater than 0")
	}

	if m.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}

	return nil
}

func (m *MetricReading) IsEmpty() bool {
	return m.ID == uuid.Nil || m.MetricID == uuid.Nil || m.Timestamp.IsZero()
}

type MetricReadingRepository interface {
	CreateMetricReading(metricReading MetricReading) (MetricReading, error)
	GetLastMetricReading(metricID uuid.UUID) (MetricReading, error)
	GetAllReadingsByMetricID(metricID uuid.UUID) ([]MetricReading, error)
}
