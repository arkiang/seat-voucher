package usecase

import (
	"bookcabin-voucher/internal/dto"
	"bookcabin-voucher/internal/model"
	"bookcabin-voucher/internal/utils"
	mockRep "bookcabin-voucher/mocks/repository"
	mockSvc "bookcabin-voucher/mocks/service"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestGenerateAndAssignSeats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRep.NewMockFlightRepository(ctrl)
	mockGen := mockSvc.NewMockSeatAllocator(ctrl)
	uc := NewFlightUsecase(mockRepo, mockGen)

	req := dto.GenerateRequest{
		CrewName:      "ApArki",
		CrewID:        "270123",
		FlightNumber:  "JT692",
		Date:          "26-07-25",
		Aircraft:      "Airbus 320",
		SeatsToChange: make([]string, 0),
	}

	seats := []model.FlightSeatAssignment{{
		Seat: "3B",
	}, {Seat: "7C"}, {Seat: "14D"}}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	tx := db.Begin()
	mockRepo.EXPECT().BeginTx().Return(tx)
	mockRepo.EXPECT().CountByFlightAndDateTx(gomock.Any(), "JT692", "26-07-25").Return(int64(0))
	mockGen.EXPECT().GenerateSeats(model.Airbus320, 3, make([]string, 0)).Return([]string{"3B", "7C", "14D"}, nil)
	mockRepo.EXPECT().BulkCreateSeatAssignmentsTx(gomock.Any(), gomock.Any()).Return(nil)
	mockRepo.EXPECT().GetByFilter(dto.FlightFilter{
		FlightNumber: "JT692", Date: "26-07-25",
	}).Return([]model.FlightAssignment{{SeatAssignments: seats}}, nil)

	expectResp := &model.FlightAssignment{
		CrewName:        req.CrewName,
		CrewID:          req.CrewID,
		FlightNumber:    req.FlightNumber,
		FlightDate:      req.Date,
		AircraftType:    req.Aircraft,
		SeatAssignments: seats,
	}
	mockRepo.EXPECT().CreateTx(gomock.Any(), gomock.Any()).Return(expectResp, nil)

	result, err := uc.GenerateAndAssignSeats(req)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(result.SeatAssignments))

	expectedSeats := utils.ExtractSeats(result.SeatAssignments)
	assert.Equal(t, []string{"3B", "7C", "14D"}, expectedSeats)
}

func TestGenerateAndChangeSeats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRep.NewMockFlightRepository(ctrl)
	mockGen := mockSvc.NewMockSeatAllocator(ctrl)
	uc := NewFlightUsecase(mockRepo, mockGen)

	req := dto.GenerateRequest{
		CrewName:      "ApArki",
		CrewID:        "270123",
		FlightNumber:  "JT692",
		Date:          "26-07-25",
		Aircraft:      "Airbus 320",
		SeatsToChange: []string{"14D"},
	}

	seatsToChange := []model.FlightSeatAssignment{{Seat: "14D"}}
	seats := []model.FlightSeatAssignment{{
		Seat: "3B",
	}, {Seat: "7C"}, {Seat: "12A"}}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	tx := db.Begin()
	mockRepo.EXPECT().BeginTx().Return(tx)
	mockRepo.EXPECT().CountByFlightAndDateTx(gomock.Any(), "JT692", "26-07-25").Return(int64(1))
	mockGen.EXPECT().GenerateSeats(model.Airbus320, 1, []string{"14D"}).Return([]string{"12A"}, nil)
	mockRepo.EXPECT().GetByFilterTx(gomock.Any(), dto.FlightFilter{
		FlightNumber: "JT692", Date: "26-07-25", Seats: []string{"14D"},
	}).Return([]model.FlightAssignment{{SeatAssignments: seatsToChange}}, nil)
	mockRepo.EXPECT().DeleteSeatsByFilterTx(gomock.Any(), gomock.Any()).Return(int64(1), nil)
	mockRepo.EXPECT().BulkCreateSeatAssignmentsTx(gomock.Any(), gomock.Any()).Return(nil)
	mockRepo.EXPECT().GetByFilter(dto.FlightFilter{
		FlightNumber: "JT692", Date: "26-07-25",
	}).Return([]model.FlightAssignment{{SeatAssignments: seats}}, nil)

	result, err := uc.GenerateAndAssignSeats(req)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(result.SeatAssignments))

	expectedSeats := utils.ExtractSeats(result.SeatAssignments)
	assert.Equal(t, []string{"3B", "7C", "12A"}, expectedSeats)
}

func TestGenerateAndAssignSeats_FlightExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockRep.NewMockFlightRepository(ctrl)
	gen := mockSvc.NewMockSeatAllocator(ctrl)
	uc := NewFlightUsecase(repo, gen)

	req := dto.GenerateRequest{
		CrewName:      "ApArki",
		CrewID:        "270123",
		FlightNumber:  "JT692",
		Date:          "26-07-25",
		Aircraft:      "Airbus 320",
		SeatsToChange: make([]string, 0),
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	tx := db.Begin()
	repo.EXPECT().BeginTx().Return(tx)
	repo.EXPECT().CountByFlightAndDateTx(gomock.Any(), "JT692", "26-07-25").Return(int64(1))

	result, err := uc.GenerateAndAssignSeats(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "assignment for this flight and date already exists and no seats to change")
}

func TestGenerateAndAssignSeats_SeatGenerationFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockRep.NewMockFlightRepository(ctrl)
	gen := mockSvc.NewMockSeatAllocator(ctrl)
	uc := NewFlightUsecase(repo, gen)

	req := dto.GenerateRequest{
		CrewName:      "ApArki",
		CrewID:        "270123",
		FlightNumber:  "JT692",
		Date:          "26-07-25",
		Aircraft:      "Airbus 320",
		SeatsToChange: make([]string, 0),
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	tx := db.Begin()
	repo.EXPECT().BeginTx().Return(tx)
	repo.EXPECT().CountByFlightAndDateTx(gomock.Any(), "JT692", "26-07-25").Return(int64(0))
	gen.EXPECT().GenerateSeats(model.Airbus320, 3, make([]string, 0)).Return(nil, errors.New("unknown aircraft"))

	result, err := uc.GenerateAndAssignSeats(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to generate seats")
}

func TestGenerateAndAssignSeats_DBCreateFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockRep.NewMockFlightRepository(ctrl)
	gen := mockSvc.NewMockSeatAllocator(ctrl)
	uc := NewFlightUsecase(repo, gen)

	req := dto.GenerateRequest{
		CrewName:     "ApArki",
		CrewID:       "270123",
		FlightNumber: "JT692",
		Date:         "26-07-25",
		Aircraft:     "Airbus 320",
	}

	seats := []string{"3B", "7C", "14D"}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	tx := db.Begin()
	repo.EXPECT().BeginTx().Return(tx)
	repo.EXPECT().CountByFlightAndDateTx(gomock.Any(), "JT692", "26-07-25").Return(int64(0))
	gen.EXPECT().GenerateSeats(model.Airbus320, 3, make([]string, 0)).Return(seats, nil)
	repo.EXPECT().CreateTx(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to create in DB"))

	result, err := uc.GenerateAndAssignSeats(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create in DB")
}
