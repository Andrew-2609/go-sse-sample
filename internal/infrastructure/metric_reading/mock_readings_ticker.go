package metric_reading

import (
	"log"
	"math/rand/v2"
	"time"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/entity"
	"github.com/Andrew-2609/go-sse-sample/internal/domain/enum"
	"github.com/Andrew-2609/go-sse-sample/internal/presentation/dto"
	"github.com/Andrew-2609/go-sse-sample/pkg/sse"
	"github.com/google/uuid"
)

type MockReadingsTicker struct {
	metricRepository        entity.MetricRepository
	metricReadingRepository entity.MetricReadingRepository
	interval                time.Duration
	sseHub                  *sse.SSEHub
	stop                    chan struct{}
}

func NewMockReadingsTicker(metricRepository entity.MetricRepository, metricReadingRepository entity.MetricReadingRepository, interval time.Duration) *MockReadingsTicker {
	return &MockReadingsTicker{
		metricRepository:        metricRepository,
		metricReadingRepository: metricReadingRepository,
		interval:                interval,
		sseHub:                  sse.GetSSEHub(),
		stop:                    make(chan struct{}),
	}
}

func (t *MockReadingsTicker) Start() {
	go func() {
		ticker := time.NewTicker(t.interval)
		defer func() {
			ticker.Stop()
			close(t.stop)
		}()

		for {
			select {
			case <-ticker.C:
				metrics, err := t.metricRepository.GetAllMetrics()
				if err != nil {
					log.Printf("error getting all metrics: %s", err)
					continue
				}

				for _, metric := range metrics {
					lastReadingValue, err := t.metricReadingRepository.GetLastMetricReadingValue(metric.ID)

					if err != nil {
						log.Printf("error getting last metric reading for metric %s: %s", metric.ID, err)
						continue
					}

					newMetricReadingID, err := uuid.NewV7()

					randomIncrease := rand.Float64() * 10
					newReadingValue := lastReadingValue + randomIncrease

					newMetricReading, err := entity.NewMetricReading(newMetricReadingID, metric.ID, newReadingValue, nil)
					if err != nil {
						log.Printf("error creating new metric reading for metric %s: %s", metric.ID, err)
						continue
					}

					_, err = t.metricReadingRepository.CreateMetricReading(newMetricReading)

					if err != nil {
						log.Printf("error creating new metric reading for metric %s: %s", metric.ID, err)
						continue
					}

					newMetricReadingResponse := dto.NewCreateMetricReadingResponseDTO(newMetricReading)

					t.sseHub.Broadcast <- sse.NewEvent(enum.EventTypeMetricReadingCreated, newMetricReadingResponse)
				}
			case <-t.stop:
				log.Println("stopping mock readings ticker")
				return
			}
		}
	}()
}

func (t *MockReadingsTicker) Stop() {
	t.stop <- struct{}{}
}
