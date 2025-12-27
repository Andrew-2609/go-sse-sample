package handler

import (
	"github.com/Andrew-2609/go-sse-sample/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

func SetupMetricHTTPGinHandler(router *gin.Engine, metricController *controller.MetricController) {
	metricsGroup := router.Group("/metrics")
	metricsGroup.POST("", metricController.CreateMetric)
	metricsGroup.GET("/:id", metricController.GetMetricByID)
}
