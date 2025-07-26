package service

import (
	"bookcabin-voucher/internal/model"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type SeatGenerator struct {
	layouts map[model.AircraftType]model.AircraftLayout
}

func NewSeatAllocator(path string) *SeatGenerator {
	file, err := os.ReadFile(path)
	if err != nil {
		panic("Failed to read seat layout file: " + err.Error())
	}
	var layouts map[model.AircraftType]model.AircraftLayout
	if err := json.Unmarshal(file, &layouts); err != nil {
		panic("Invalid JSON in seat layout file: " + err.Error())
	}
	return &SeatGenerator{layouts: layouts}
}

func (s *SeatGenerator) GenerateSeats(aircraft model.AircraftType, count int) ([]string, error) {
	layout, ok := s.layouts[aircraft]
	if !ok {
		return nil, fmt.Errorf("unknown aircraft")
	}
	result := make([]string, 0, count)
	seen := make(map[string]bool)
	tries := 0
	maxTries := 1000
	for len(result) < count && tries < maxTries {
		row := rand.Intn(layout.EndRow-layout.StartRow+1) + layout.StartRow
		seat := fmt.Sprintf("%d%s", row, layout.Seats[rand.Intn(len(layout.Seats))])
		if !seen[seat] {
			result = append(result, seat)
			seen[seat] = true
		}
		tries++
	}
	if len(result) < count {
		return nil, fmt.Errorf("not enough available seats")
	}
	return result, nil
}

var _ SeatAllocator = (*SeatGenerator)(nil)
