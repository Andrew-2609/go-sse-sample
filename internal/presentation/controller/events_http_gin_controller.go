package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Andrew-2609/go-sse-sample/pkg/sse"
	"github.com/gin-gonic/gin"
)

type EventsController struct {
	sseHub *sse.SSEHub
}

func NewEventsController() *EventsController {
	return &EventsController{
		sseHub: sse.GetSSEHub(),
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

	if err := c.printDataMessage(ctx.Writer, "connected"); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if lastEventID := ctx.GetHeader("Last-Event-ID"); lastEventID != "" {
		events := c.sseHub.GetEventsAfterID(lastEventID)
		c.sendEvents(ctx.Writer, events...)
	}

	flusher.Flush()

	// any `return` triggers defer -> unregister client
	for {
		select {
		case isDisconnected := <-client.IsDisconnected():
			if isDisconnected {
				connDuration := time.Since(connStartTime)
				log.Printf("client disconnected after %d seconds", int(math.Ceil(connDuration.Seconds())))

				if err := c.printDataMessage(ctx.Writer, "disconnected"); err != nil {
					return
				}

				flusher.Flush()
				return
			}
		case event := <-client.CH():
			if err := c.sendEvents(ctx.Writer, event); err != nil {
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

func (c *EventsController) sendEvents(w io.Writer, events ...sse.Event) error {
	printLines := func(lines ...string) error {
		for _, line := range lines {
			if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
				return fmt.Errorf("error sending line: %w", err)
			}
		}

		return nil
	}

	for _, event := range events {
		if event.IsEmpty() {
			continue
		}

		eventData, err := json.Marshal(event.Data)
		if err != nil {
			return fmt.Errorf("error marshalling event data: %w", err)
		}

		lines := []string{
			fmt.Sprintf("id: %s", event.ID),
			fmt.Sprintf("event: %s", event.Type),
			fmt.Sprintf("data: %s", string(eventData)),
		}

		if err := printLines(lines...); err != nil {
			return err
		}

		// end of event (CRITICAL)
		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return fmt.Errorf("error sending event end of line: %w", err)
		}
	}

	return nil
}

func (c *EventsController) printDataMessage(w io.Writer, message string, args ...any) error {
	if _, err := fmt.Fprintf(w, "data: %s\n\n", fmt.Sprintf(message, args...)); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}
