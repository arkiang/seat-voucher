package handler

import (
	apiModel "bookcabin-voucher/internal/api/model"
	"bookcabin-voucher/internal/dto"
	"bookcabin-voucher/internal/model"
	"bookcabin-voucher/internal/validation"
	mockUc "bookcabin-voucher/mocks/usecase"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckFlightHandler_Exists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validation.RegisterValidators()

	mockUsecase := mockUc.NewMockFlightUsecase(ctrl)
	h := NewFlightHandler(mockUsecase)
	r := gin.Default()
	r.POST("/api/check", h.CheckFlight)

	mockUsecase.EXPECT().CheckFlightExists(dto.CheckFlightRequest{
		FlightNumber: "JT692",
		Date:         "26-07-25",
	}).Return(true)

	reqBody := `{"flightNumber":"JT692","date":"26-07-25"}`
	req := httptest.NewRequest(http.MethodPost, "/api/check", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"exists": true}`, resp.Body.String())
}

func TestGenerateFlightHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validation.RegisterValidators()

	mockUsecase := mockUc.NewMockFlightUsecase(ctrl)
	h := NewFlightHandler(mockUsecase)
	r := gin.Default()
	r.POST("/api/generate", h.Generate)

	reqData := dto.GenerateRequest{
		CrewName:     "ApArki",
		CrewID:       "98123",
		FlightNumber: "JT692",
		Date:         "26-07-25",
		Aircraft:     "Airbus 320",
	}

	assignment := &model.FlightAssignment{
		CrewName:        "ApArki",
		CrewID:          "98123",
		FlightNumber:    "JT692",
		FlightDate:      "26-07-25",
		AircraftType:    "Airbus 320",
		SeatAssignments: []model.FlightSeatAssignment{{Seat: "3A"}, {Seat: "5C"}, {Seat: "8F"}},
	}

	mockUsecase.EXPECT().GenerateAndAssignSeats(reqData).Return(assignment, nil)

	body, _ := json.Marshal(reqData)
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var bodyResp dto.GenerateResponse
	err := json.Unmarshal(resp.Body.Bytes(), &bodyResp)
	assert.NoError(t, err)
	assert.Equal(t, bodyResp, dto.GenerateResponse{Success: true, Seats: []string{"3A", "5C", "8F"}})
}

func TestGenerateFlightHandler_ValidationAircraftFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validation.RegisterValidators()

	mockUsecase := mockUc.NewMockFlightUsecase(ctrl)
	h := NewFlightHandler(mockUsecase)
	r := gin.Default()
	r.POST("/api/generate", h.Generate)

	// missing crewID and wrong aircraft
	reqBody := `{
		"crewName": "",
		"crewID": "",
		"flightNumber": "JT692",
		"date": "12-07-25",
		"aircraft": "INVALID"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var body apiModel.ErrorResponse
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Contains(t, body.Error, "Invalid input")
}

func TestGenerateFlightHandler_ValidationFlightNumberFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validation.RegisterValidators()

	mockUsecase := mockUc.NewMockFlightUsecase(ctrl)
	h := NewFlightHandler(mockUsecase)
	r := gin.Default()
	r.POST("/api/generate", h.Generate)

	// wrong flight number
	reqBody := `{
		"crewName": "ApArki",
		"crewID": "122511",
		"flightNumber": "ID12345",
		"date": "12-07-25",
		"aircraft": "Airbus 320"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var body apiModel.ErrorResponse
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Contains(t, body.Error, "Invalid input")
}

func TestGenerateFlightHandler_UsecaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockUc.NewMockFlightUsecase(ctrl)
	h := NewFlightHandler(mockUsecase)
	r := gin.Default()
	r.POST("/api/generate", h.Generate)

	reqData := dto.GenerateRequest{
		CrewName:     "Sarah",
		CrewID:       "98123",
		FlightNumber: "JT692",
		Date:         "12-07-25",
		Aircraft:     "Airbus 320",
	}

	mockUsecase.EXPECT().GenerateAndAssignSeats(reqData).
		Return(nil, fmt.Errorf("assignment for this flight and date already exists"))

	body, _ := json.Marshal(reqData)
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	var bodyResp apiModel.ErrorResponse
	err := json.Unmarshal(resp.Body.Bytes(), &bodyResp)
	assert.NoError(t, err)
	assert.Equal(t, bodyResp.Error, "assignment for this flight and date already exists")
}
