package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
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

var (
	depsOnce                sync.Once
	metricController        *controller.MetricController
	metricReadingController *controller.MetricReadingController
	eventsController        *controller.EventsController
	eventStore              *repository.EventStoreInMemory
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	setupDependencies()

	router := gin.Default()
	router.Use(corsMiddleware())

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

	eventStore.StopRetention()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v\n", err)
	}

	log.Println("server exiting")
}

func setupRoutes(router *gin.Engine) {
	metricsGroup := router.Group("/metrics")
	metricReadingsGroup := metricsGroup.Group("/readings")
	eventsGroup := router.Group("/events")

	metricController.SetupRoutes(metricsGroup)
	metricReadingController.SetupRoutes(metricReadingsGroup)
	eventsController.SetupRoutes(eventsGroup)
}

func setupDependencies() {
	depsOnce.Do(func() {
		inMemoryEventsTTL := 1 * time.Minute
		eventStore = repository.NewEventStoreInMemory(inMemoryEventsTTL)
		sse.InitializeSSEHub(eventStore, MAX_SSE_CLIENTS)

		metricRepository := repository.NewMetricInMemoryRepository()
		metricUseCase := use_case.NewMetricUseCase(metricRepository)
		metricController = controller.NewMetricController(metricUseCase)

		metricReadingRepository := repository.NewMetricReadingInMemoryRepository()
		metricReadingUseCase := use_case.NewMetricReadingUseCase(metricRepository, metricReadingRepository)
		metricReadingController = controller.NewMetricReadingController(metricReadingUseCase)

		eventsController = controller.NewEventsController()
	})
}

// Non-realistic CORS, for development only!
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Last-Event-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
