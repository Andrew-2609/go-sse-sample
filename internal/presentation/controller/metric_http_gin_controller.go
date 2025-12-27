package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/use_case"
	"github.com/Andrew-2609/go-sse-sample/internal/presentation/dto"
	"github.com/Andrew-2609/go-sse-sample/pkg/sse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MetricController struct {
	metricUseCase *use_case.MetricUseCase
	sseHub        *sse.SSEHub
}

func NewMetricController(metricUseCase *use_case.MetricUseCase, sseHub *sse.SSEHub) *MetricController {
	return &MetricController{
		metricUseCase: metricUseCase,
		sseHub:        sseHub,
	}
}

func (c *MetricController) SetupRoutes(metricsGroup *gin.RouterGroup) {
	metricsGroup.POST("", c.CreateMetric)
	metricsGroup.GET("/:id", c.GetMetricByID)
	metricsGroup.GET("/watch-created", c.WatchCreatedMetrics)
}

func (c *MetricController) CreateMetric(ctx *gin.Context) {
	var request dto.CreateMetricRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.metricUseCase.CreateMetric(request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *MetricController) GetMetricByID(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.metricUseCase.GetMetricByID(parsedID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *MetricController) WatchCreatedMetrics(ctx *gin.Context) {
	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "flusher not supported"})
		return
	}

	closeNotify := ctx.Writer.CloseNotify()

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	connStartTime := time.Now().UTC()

	client := sse.NewSSEClient(
		make(chan sse.Event, 8),
		connStartTime,
	)

	c.sseHub.Register <- client

	log.Printf("new client connected at %s\n", connStartTime.Format(time.RFC3339))

	defer func() {
		c.sseHub.Unregister <- client
	}()

	type connectionMessage struct {
		Message string `json:"message"`
	}

	message := connectionMessage{
		Message: fmt.Sprintf("client connected at %s", connStartTime.Format(time.RFC3339)),
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshalling message: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := fmt.Fprintf(ctx.Writer, "%s\n", string(data)); err != nil {
		log.Printf("error sending message: %v\n", err)
		return
	}

	flusher.Flush()

	// any `return` triggers defer -> unregister client
	for {
		select {
		case isDisconnected := <-client.IsDisconnected():
			if isDisconnected {
				message := connectionMessage{
					Message: fmt.Sprintf("client disconnected at %s", time.Now().UTC().Format(time.RFC3339)),
				}

				data, err := json.Marshal(message)
				if err != nil {
					log.Printf("error marshalling message: %v\n", err)
					return
				}

				if _, err := fmt.Fprintf(ctx.Writer, "%s\n", string(data)); err != nil {
					log.Printf("error sending message: %v\n", err)
					return
				}

				flusher.Flush()
				return
			}
		case event := <-client.CH():
			if event.IsEmpty() {
				continue
			}

			data, err := json.Marshal(event)
			if err != nil {
				log.Printf("error marshalling event: %v\n", err)
				return
			}

			if _, err := fmt.Fprintf(ctx.Writer, "%s\n", string(data)); err != nil {
				log.Printf("error sending event: %v\n", err)
				return
			}

			flusher.Flush()
		case <-closeNotify:
			connDuration := time.Since(connStartTime)
			log.Printf("client disconnected after %d seconds", int(math.Ceil(connDuration.Seconds())))
			return
		}
	}
}
