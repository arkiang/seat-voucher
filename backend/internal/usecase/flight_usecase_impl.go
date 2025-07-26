package usecase

import (
	"bookcabin-voucher/internal/dto"
	"bookcabin-voucher/internal/model"
	"bookcabin-voucher/internal/repository"
	"bookcabin-voucher/internal/service"
	"fmt"
	"log"
	"strings"
)

type flightUsecaseImpl struct {
	repo    repository.FlightRepository
	seatGen service.SeatAllocator
}

func NewFlightUsecase(repo repository.FlightRepository, seatGen service.SeatAllocator) FlightUsecase {
	return &flightUsecaseImpl{
		repo:    repo,
		seatGen: seatGen,
	}
}

func (u *flightUsecaseImpl) CheckFlightExists(request dto.CheckFlightRequest) bool {
	return u.repo.CountByFlightAndDate(request.FlightNumber, request.Date) > 0
}

func (u *flightUsecaseImpl) GenerateAndAssignSeats(request dto.GenerateRequest) (*model.FlightAssignment, error) {
	if u.CheckFlightExists(dto.CheckFlightRequest{FlightNumber: request.FlightNumber, Date: request.Date}) {
		log.Printf("[Usecase] Flight assignment already exists: %s on %s", request.FlightNumber, request.Date)
		return nil, fmt.Errorf("assignment for this flight and date already exists")
	}

	seats, err := u.seatGen.GenerateSeats(request.Aircraft, 3)
	if err != nil {
		log.Printf("[Usecase] Seat generation failed for %s: %v", request.Aircraft, err)
		return nil, fmt.Errorf("failed to generate seats: %w", err)
	}

	assignment := &model.FlightAssignment{
		CrewName:     request.CrewName,
		CrewID:       request.CrewID,
		FlightNumber: request.FlightNumber,
		FlightDate:   request.Date,
		AircraftType: request.Aircraft,
		Seats:        strings.Join(seats, ","),
	}

	result, err := u.repo.Create(assignment)
	if err != nil {
		log.Printf("[Usecase] Failed to persist assignment for %s: %v", request.FlightNumber, err)
		return nil, fmt.Errorf("failed to create in DB: %w", err)
	}

	return result, nil
}
