package usecase

import (
	"bookcabin-voucher/internal/dto"
	"bookcabin-voucher/internal/model"
)

type FlightUsecase interface {
	CheckFlightExists(request dto.CheckFlightRequest) bool
	GenerateAndAssignSeats(request dto.GenerateRequest) (*model.FlightAssignment, error)
}
