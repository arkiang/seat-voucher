package repository

import "bookcabin-voucher/internal/model"

type FlightRepository interface {
	CountByFlightAndDate(flightNumber, date string) int64
	Create(assignment *model.FlightAssignment) (*model.FlightAssignment, error)
}
