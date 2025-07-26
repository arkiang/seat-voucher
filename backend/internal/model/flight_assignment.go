package model

import "time"

type AircraftType string

const (
	ATR          AircraftType = "ATR"
	Airbus320    AircraftType = "Airbus 320"
	Boeing737Max AircraftType = "Boeing 737 Max"
)

type FlightAssignment struct {
	ID           uint         `gorm:"primaryKey"`
	CrewName     string       `gorm:"type:varchar(100);not null"`
	CrewID       string       `gorm:"type:varchar(50);not null"`
	FlightNumber string       `gorm:"type:varchar(20);not null"`
	FlightDate   string       `gorm:"type:date;not null"`
	AircraftType AircraftType `gorm:"type:varchar(50);not null"`
	Seats        string       `gorm:"type:text;not null"` // comma-separated

	CreatedAt time.Time `gorm:"autoCreateTime"`
}
