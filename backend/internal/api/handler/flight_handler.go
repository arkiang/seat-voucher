package handler

import (
	"bookcabin-voucher/internal/api/model"
	"bookcabin-voucher/internal/dto"
	serviceModel "bookcabin-voucher/internal/model"
	"bookcabin-voucher/internal/usecase"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
		Seats:   splitSeats(assignment.SeatAssignments),
	})
}

func splitSeats(seats []serviceModel.FlightSeatAssignment) []string {
	if len(seats) == 0 {
		return []string{}
	}

	var result []string
	for _, p := range seats {
		result = append(result, p.Seat)
	}
	return result
}
