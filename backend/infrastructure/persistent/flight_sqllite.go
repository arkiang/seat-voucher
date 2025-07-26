package persistent

import (
	"bookcabin-voucher/internal/model"
	"bookcabin-voucher/internal/repository"
	"gorm.io/gorm"
)

type flightRepository struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) repository.FlightRepository {
	return &flightRepository{db: db}
}

func (r *flightRepository) CountByFlightAndDate(flightNumber, date string) int64 {
	var count int64
	r.db.Model(&model.FlightAssignment{}).
		Where("flight_number = ? AND flight_date = ?", flightNumber, date).
		Count(&count)
	return count
}

func (r *flightRepository) Create(assignment *model.FlightAssignment) (*model.FlightAssignment, error) {
	err := r.db.Create(assignment).Error
	return assignment, err
}
