package model

type AircraftLayout struct {
	StartRow int      `json:"startRow"`
	EndRow   int      `json:"endRow"`
	Seats    []string `json:"seats"`
}
