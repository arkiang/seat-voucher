package service

import (
	"bookcabin-voucher/config"
	"bookcabin-voucher/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupTestLayout(t *testing.T) *SeatGenerator {
	cfg := config.LoadConfig()

	layoutPath := cfg.SeatLayoutPath
	assert.NotEmpty(t, layoutPath)

	gen := NewSeatAllocator(layoutPath)
	return gen
}

func TestGenerateSeats_Success(t *testing.T) {
	gen := setupTestLayout(t)
	seats, err := gen.GenerateSeats(model.Airbus320, 3, make([]string, 0))

	assert.NoError(t, err)
	assert.Len(t, seats, 3)
}

func TestGenerateSeats_UnknownAircraft(t *testing.T) {
	gen := setupTestLayout(t)
	seats, err := gen.GenerateSeats("some-unknown", 3, make([]string, 0))

	assert.Error(t, err)
	assert.Nil(t, seats)
	assert.Contains(t, err.Error(), "unknown aircraft")
}

func TestGenerateSeats_InsufficientSeats(t *testing.T) {
	gen := setupTestLayout(t)
	seats, err := gen.GenerateSeats(model.Airbus320, 50000000, make([]string, 0))

	assert.Error(t, err)
	assert.Nil(t, seats)
	assert.Contains(t, err.Error(), "not enough available seats")
}
