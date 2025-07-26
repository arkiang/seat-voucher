package usecase

import (
	"bookcabin-voucher/internal/dto"
	"bookcabin-voucher/internal/model"
	mockRep "bookcabin-voucher/mocks/repository"
	mockSvc "bookcabin-voucher/mocks/service"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGenerateAndAssignSeats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRep.NewMockFlightRepository(ctrl)
	mockGen := mockSvc.NewMockSeatAllocator(ctrl)
	uc := NewFlightUsecase(mockRepo, mockGen)

	req := dto.GenerateRequest{
		CrewName:     "ApArki",
		CrewID:       "270123",
		FlightNumber: "JT692",
		Date:         "26-07-25",
		Aircraft:     "Airbus 320",
	}

	mockRepo.EXPECT().CountByFlightAndDate("JT692", "26-07-25").Return(int64(0))
	mockGen.EXPECT().GenerateSeats(model.Airbus320, 3).Return([]string{"3B", "7C", "14D"}, nil)
	mockRepo.EXPECT().Create(gomock.Any()).Return(&model.FlightAssignment{Seats: "3B,7C,14D"}, nil)

	result, err := uc.GenerateAndAssignSeats(req)

	assert.NoError(t, err)
	assert.Equal(t, "3B,7C,14D", result.Seats)
}

func TestGenerateAndAssignSeats_FlightExists(t *testing.T) {
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

	repo.EXPECT().CountByFlightAndDate("JT692", "26-07-25").Return(int64(1))

	result, err := uc.GenerateAndAssignSeats(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "assignment for this flight and date already exists")
}

func TestGenerateAndAssignSeats_SeatGenerationFailed(t *testing.T) {
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

	repo.EXPECT().CountByFlightAndDate("JT692", "26-07-25").Return(int64(0))
	gen.EXPECT().GenerateSeats(model.Airbus320, 3).Return(nil, errors.New("unknown aircraft"))

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

	repo.EXPECT().CountByFlightAndDate("JT692", "26-07-25").Return(int64(0))
	gen.EXPECT().GenerateSeats(model.Airbus320, 3).Return(seats, nil)
	repo.EXPECT().Create(gomock.Any()).Return(nil, errors.New("failed to create in DB"))

	result, err := uc.GenerateAndAssignSeats(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create in DB")
}
