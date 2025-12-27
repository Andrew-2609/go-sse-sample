package controller

import (
	"net/http"
	"time"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/use_case"
	"github.com/Andrew-2609/go-sse-sample/internal/presentation/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MetricReadingController struct {
	metricReadingUseCase *use_case.MetricReadingUseCase
}

func NewMetricReadingController(metricReadingUseCase *use_case.MetricReadingUseCase) *MetricReadingController {
	return &MetricReadingController{
		metricReadingUseCase: metricReadingUseCase,
	}
}

func (c *MetricReadingController) SetupRoutes(metricReadingsGroup *gin.RouterGroup) {
	metricReadingsGroup.POST("", c.CreateMetricReading)
}

func (c *MetricReadingController) CreateMetricReading(ctx *gin.Context) {
	var request dto.CreateMetricReadingRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := uuid.Parse(request.MetricID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.HasTimestamp() {
		if _, err := time.Parse(time.RFC3339, *request.Timestamp); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	response, err := c.metricReadingUseCase.CreateMetricReading(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
