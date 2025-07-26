package service

import "bookcabin-voucher/internal/model"

type SeatAllocator interface {
	GenerateSeats(aircraft model.AircraftType, count int) ([]string, error)
}
