package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/use_case"
	"github.com/Andrew-2609/go-sse-sample/internal/presentation/controller"
	"github.com/Andrew-2609/go-sse-sample/internal/repository"
	"github.com/Andrew-2609/go-sse-sample/pkg/sse"
	"github.com/gin-gonic/gin"
)

const (
	MAX_SSE_CLIENTS = 10_000
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	setupRoutes(router)

	srv := &http.Server{
		Addr:    ":8089",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("server started on port 8089")

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v\n", err)
	}

	log.Println("server exiting")
}

func setupRoutes(router *gin.Engine) {
	metricController, metricReadingController, eventsController := makeControllers()

	metricsGroup := router.Group("/metrics")
	metricReadingsGroup := metricsGroup.Group("/readings")
	eventsGroup := router.Group("/events")

	metricController.SetupRoutes(metricsGroup)
	metricReadingController.SetupRoutes(metricReadingsGroup)
	eventsController.SetupRoutes(eventsGroup)
}

func makeControllers() (
	*controller.MetricController,
	*controller.MetricReadingController,
	*controller.EventsController,
) {
	eventStore := repository.NewEventStoreInMemory()
	sseHub := sse.NewSSEHub(eventStore, MAX_SSE_CLIENTS)

	metricRepository := repository.NewMetricInMemoryRepository()
	metricUseCase := use_case.NewMetricUseCase(metricRepository, sseHub)
	metricController := controller.NewMetricController(metricUseCase)

	metricReadingRepository := repository.NewMetricReadingInMemoryRepository()
	metricReadingUseCase := use_case.NewMetricReadingUseCase(metricRepository, metricReadingRepository)
	metricReadingController := controller.NewMetricReadingController(metricReadingUseCase)

	eventsController := controller.NewEventsController(sseHub)

	return metricController, metricReadingController, eventsController
}
