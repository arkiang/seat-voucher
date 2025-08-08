package migration

import (
	"bookcabin-voucher/internal/model"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.FlightAssignment{},
		&model.FlightSeatAssignment{},
	)

	if err != nil {
		log.Fatalf("Failed to run AutoMigrate: %v", err)
	}

	// Enforce unique constraint on (flight_number, date)
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_flight_number_date ON flight_assignments(flight_number, flight_date)")
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_flight_assignment_id_seat ON flight_seat_assignments(flight_assignment_id, seat)")

	log.Println("Database migration completed.")
}
