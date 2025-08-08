package usecase

import (
	"bookcabin-voucher/internal/dto"
	"bookcabin-voucher/internal/model"
	"bookcabin-voucher/internal/repository"
	"bookcabin-voucher/internal/service"
	"bookcabin-voucher/internal/utils"
	"fmt"
	"log"
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
	tx := u.repo.BeginTx()
	count := u.repo.CountByFlightAndDateTx(tx, request.FlightNumber, request.Date)

	//if not exist, create new
	if count == 0 {
		seats, err := u.seatGen.GenerateSeats(request.Aircraft, 3, make([]string, 0))
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
		}

		//create assignment
		assignment, err = u.repo.CreateTx(tx, assignment)
		if err != nil {
			tx.Rollback()
			log.Printf("[Usecase] Failed to persist assignment for %s: %v", request.FlightNumber, err)
			return nil, fmt.Errorf("failed to create assignment in DB: %w", err)
		}

		var seatAssignments []model.FlightSeatAssignment
		for _, seat := range seats {
			seatAssignments = append(seatAssignments, model.FlightSeatAssignment{
				FlightAssignmentID: assignment.ID,
				Seat:               seat,
			})
		}

		//create seat assignment
		err = u.repo.BulkCreateSeatAssignmentsTx(tx, seatAssignments)
		if err != nil {
			tx.Rollback()
			log.Printf("[Usecase] Failed to create seat assignments for %s: %v", request.FlightNumber, err)
			return nil, fmt.Errorf("failed to create seat assignments: %w", err)
		}
	} else { //if exist will use update instead
		seatsToChangeCount := len(request.SeatsToChange)
		if seatsToChangeCount == 0 {
			tx.Rollback()
			log.Printf("[Usecase] Flight assignment already exists: %s on %s", request.FlightNumber, request.Date)
			return nil, fmt.Errorf("assignment for this flight and date already exists and no seats to change")
		}

		filter := dto.FlightFilter{
			FlightNumber: request.FlightNumber,
			Date:         request.Date,
			Seats:        request.SeatsToChange,
		}

		// Find the related assignment for this flight
		assignments, err := u.repo.GetByFilterTx(tx, filter)
		if err != nil || len(assignments) == 0 {
			tx.Rollback()
			log.Printf("[Usecase] No assignment found after deletion or error: %v", err)
			return nil, fmt.Errorf("no matching assignment found after seat deletion")
		}

		//generate new seats assignment
		seatsToChange := utils.ExtractSeats(assignments[0].SeatAssignments)
		seats, err := u.seatGen.GenerateSeats(request.Aircraft, seatsToChangeCount, seatsToChange)
		if err != nil {
			log.Printf("[Usecase] Seat generation failed for %s: %v", request.Aircraft, err)
			return nil, fmt.Errorf("failed to generate seats: %w", err)
		}

		// Delete existing seats
		_, err = u.repo.DeleteSeatsByFilterTx(tx, filter)
		if err != nil {
			tx.Rollback()
			log.Printf("[Usecase] Failed to delete existing seats: %v", err)
			return nil, fmt.Errorf("failed to delete seats: %w", err)
		}

		var seatAssignments []model.FlightSeatAssignment
		for _, seat := range seats {
			seatAssignments = append(seatAssignments, model.FlightSeatAssignment{
				FlightAssignmentID: assignments[0].ID,
				Seat:               seat,
			})
		}

		//insert new seats
		err = u.repo.BulkCreateSeatAssignmentsTx(tx, seatAssignments)
		if err != nil {
			tx.Rollback()
			log.Printf("[Usecase] Failed to re-create seats: %v", err)
			return nil, fmt.Errorf("failed to re-create seat assignments: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	currentFilter := dto.FlightFilter{
		FlightNumber: request.FlightNumber,
		Date:         request.Date,
	}

	// find the updated data, a guarantee will be there
	assignments, _ := u.repo.GetByFilter(currentFilter)

	return &assignments[0], nil
}
