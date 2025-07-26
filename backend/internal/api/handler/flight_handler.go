package handler

import (
	"bookcabin-voucher/internal/api/model"
	"bookcabin-voucher/internal/dto"
	"bookcabin-voucher/internal/usecase"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type FlightHandler struct {
	Usecase usecase.FlightUsecase
}

func NewFlightHandler(u usecase.FlightUsecase) *FlightHandler {
	return &FlightHandler{Usecase: u}
}

func (h *FlightHandler) CheckFlight(c *gin.Context) {
	var req dto.CheckFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CheckFlight] Validation failed: %v", err)

		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Invalid input: " + err.Error(),
		})
		return
	}
	exists := h.Usecase.CheckFlightExists(req)
	c.JSON(http.StatusOK, dto.CheckFlightResponse{Exists: exists})
}

func (h *FlightHandler) Generate(c *gin.Context) {
	var req dto.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Generate] Validation failed: %v", err)

		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Invalid input: " + err.Error(),
		})
		return
	}
	assignment, err := h.Usecase.GenerateAndAssignSeats(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.GenerateResponse{
		Success: true,
		Seats:   splitSeats(assignment.Seats),
	})
}

func splitSeats(seats string) []string {
	if seats == "" {
		return []string{}
	}
	return strings.Split(seats, ",")
}
