package validation

import (
	"bookcabin-voucher/internal/model"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("flight_number", FlightNumberValidator)
		v.RegisterValidation("aircraft_enum", AircraftEnumValidator)
	}
}

// AircraftEnumValidator checks if Aircraft is one of the allowed values
func AircraftEnumValidator(fl validator.FieldLevel) bool {
	aircraft := model.AircraftType(fl.Field().String())
	switch aircraft {
	case model.ATR, model.Airbus320, model.Boeing737Max:
		return true
	default:
		return false
	}
}

var flightNumberRegex = regexp.MustCompile(`^[A-Z]{2}\d{1,4}$`)

// FlightNumberValidator checks if flightNumber following a correct pattern
func FlightNumberValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return flightNumberRegex.MatchString(value)
}
