package persistent

import (
	"bookcabin-voucher/internal/dto"
	"bookcabin-voucher/internal/model"
	"bookcabin-voucher/internal/repository"
	"fmt"
	"gorm.io/gorm"
)

type flightRepository struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) repository.FlightRepository {
	return &flightRepository{db: db}
}

func (r *flightRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *flightRepository) CountByFlightAndDate(flightNumber, date string) int64 {
	return r.CountByFlightAndDateTx(r.db, flightNumber, date)
}

func (r *flightRepository) CountByFlightAndDateTx(tx *gorm.DB, flightNumber, date string) int64 {
	var count int64
	tx.Model(&model.FlightAssignment{}).
		Where("flight_number = ? AND flight_date = ?", flightNumber, date).
		Count(&count)
	return count
}

func (r *flightRepository) GetByFilter(filter dto.FlightFilter) ([]model.FlightAssignment, error) {
	return r.GetByFilterTx(r.db, filter)
}

func (r *flightRepository) GetByFilterTx(tx *gorm.DB, filter dto.FlightFilter) ([]model.FlightAssignment, error) {
	var assignments []model.FlightAssignment

	query := tx.
		Model(&model.FlightAssignment{}).
		Preload("SeatAssignments").
		Joins("JOIN flight_seat_assignments ON flight_assignments.id = flight_seat_assignments.flight_assignment_id").
		Where("flight_assignments.flight_number = ? AND flight_assignments.flight_date = ?", filter.FlightNumber, filter.Date)

	if len(filter.Seats) > 0 {
		query = query.Where("flight_seat_assignments.seat IN ?", filter.Seats)
	}

	err := query.
		Group("flight_assignments.id").
		Order("flight_seat_assignments.created_at DESC").
		Find(&assignments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query flight assignments in transaction: %w", err)
	}

	return assignments, nil
}

func (r *flightRepository) DeleteSeatsByFilterTx(tx *gorm.DB, filter dto.FlightFilter) (int64, error) {
	var assignment model.FlightAssignment
	if err := tx.Where("flight_number = ? AND flight_date = ?", filter.FlightNumber, filter.Date).
		First(&assignment).Error; err != nil {
		return 0, err
	}

	result := tx.Where("flight_assignment_id = ? AND seat IN ?", assignment.ID, filter.Seats).
		Delete(&model.FlightSeatAssignment{})

	return result.RowsAffected, result.Error
}

func (r *flightRepository) BulkCreateSeatAssignmentsTx(tx *gorm.DB, seats []model.FlightSeatAssignment) error {
	if len(seats) == 0 {
		return nil
	}
	return tx.Create(&seats).Error
}

func (r *flightRepository) CreateTx(tx *gorm.DB, assignment *model.FlightAssignment) (*model.FlightAssignment, error) {
	if err := tx.Create(assignment).Error; err != nil {
		return nil, err
	}
	return assignment, nil
}
