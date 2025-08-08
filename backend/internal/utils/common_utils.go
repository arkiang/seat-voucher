package utils

import "bookcabin-voucher/internal/model"

func ExtractSeats(assignments []model.FlightSeatAssignment) []string {
	seats := make([]string, 0, len(assignments))
	for _, a := range assignments {
		seats = append(seats, a.Seat)
	}
	return seats
}
