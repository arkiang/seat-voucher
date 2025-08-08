package repository

import (
	"bookcabin-voucher/internal/dto"
	"bookcabin-voucher/internal/model"
	"gorm.io/gorm"
)

type FlightRepository interface {
	BeginTx() *gorm.DB

	CountByFlightAndDate(flightNumber, date string) int64
	CountByFlightAndDateTx(tx *gorm.DB, flightNumber, date string) int64
	GetByFilter(filter dto.FlightFilter) ([]model.FlightAssignment, error)
	GetByFilterTx(tx *gorm.DB, filter dto.FlightFilter) ([]model.FlightAssignment, error)

	CreateTx(tx *gorm.DB, assignment *model.FlightAssignment) (*model.FlightAssignment, error)
	DeleteSeatsByFilterTx(tx *gorm.DB, filter dto.FlightFilter) (int64, error)
	BulkCreateSeatAssignmentsTx(tx *gorm.DB, seats []model.FlightSeatAssignment) error
}
