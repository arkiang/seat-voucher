package dto

import "bookcabin-voucher/internal/model"

type CheckFlightRequest struct {
	FlightNumber string `json:"flightNumber" binding:"required,flight_number"`
	Date         string `json:"date" binding:"required,datetime=02-01-06"`
}

type CheckFlightResponse struct {
	Exists bool `json:"exists"`
}

type GenerateRequest struct {
	CrewName     string             `json:"name" binding:"required"`
	CrewID       string             `json:"id" binding:"required"`
	FlightNumber string             `json:"flightNumber" binding:"required,flight_number"`
	Date         string             `json:"date" binding:"required,datetime=02-01-06"` //DD-MM-YY
	Aircraft     model.AircraftType `json:"aircraft" binding:"required,aircraft_enum"`
}

type GenerateResponse struct {
	Success bool     `json:"success"`
	Seats   []string `json:"seats"`
}
