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
	"github.com/Andrew-2609/go-sse-sample/internal/presentation/handler"
	"github.com/Andrew-2609/go-sse-sample/internal/repository"
	"github.com/gin-gonic/gin"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v\n", err)
	}

	log.Println("server exiting")
}

func setupRoutes(router *gin.Engine) {
	metricController := makeMetricController()
	handler.SetupMetricHTTPGinHandler(router, metricController)
}

func makeMetricController() *controller.MetricController {
	metricRepository := repository.NewMetricInMemoryRepository()
	metricUseCase := use_case.NewMetricUseCase(metricRepository)
	metricController := controller.NewMetricController(metricUseCase)
	return metricController
}
