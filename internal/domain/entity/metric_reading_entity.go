package entity

import "time"

type MetricReading struct {
	ID        string
	MetricID  string
	Value     float64
	Timestamp time.Time
}

type MetricReadingRepository interface {
	CreateMetricReading(metricReading MetricReading) error
}
