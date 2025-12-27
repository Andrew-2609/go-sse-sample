package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Andrew-2609/go-sse-sample/internal/presentation/dto"
	"github.com/Andrew-2609/go-sse-sample/pkg/sse"
	"github.com/gin-gonic/gin"
)

type EventsController struct {
	sseHub *sse.SSEHub
}

func NewEventsController(sseHub *sse.SSEHub) *EventsController {
	return &EventsController{
		sseHub: sseHub,
	}
}

func (c *EventsController) SetupRoutes(eventsGroup *gin.RouterGroup) {
	eventsGroup.GET("/watch", c.WatchEvents)
}

func (c *EventsController) WatchEvents(ctx *gin.Context) {
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

	message := dto.NewConnectionMessageDTO("client connected at %s", connStartTime.Format(time.RFC3339))

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
				message := dto.NewConnectionMessageDTO("client disconnected at %s", time.Now().UTC().Format(time.RFC3339))

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
