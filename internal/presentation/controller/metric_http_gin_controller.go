package controller

import (
	"net/http"

	"github.com/Andrew-2609/go-sse-sample/internal/domain/use_case"
	"github.com/Andrew-2609/go-sse-sample/internal/presentation/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MetricController struct {
	metricUseCase *use_case.MetricUseCase
}

func NewMetricController(metricUseCase *use_case.MetricUseCase) *MetricController {
	return &MetricController{
		metricUseCase: metricUseCase,
	}
}

func (c *MetricController) SetupRoutes(metricsGroup *gin.RouterGroup) {
	metricsGroup.POST("", c.CreateMetric)
	metricsGroup.GET("/:id", c.GetMetricByID)
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
