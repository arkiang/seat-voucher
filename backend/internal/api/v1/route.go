package http

import (
	"bookcabin-voucher/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, flightHandler *handler.FlightHandler) {
	r.POST("/api/check", flightHandler.CheckFlight)
	r.POST("/api/generate", flightHandler.Generate)
}
