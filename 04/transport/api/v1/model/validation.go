package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterStructValidation(validateBirthDate, BirthDate{})
}

func validateBirthDate(sl validator.StructLevel) {
	birth := sl.Current().Interface().(BirthDate)

	maxDay := getDaysInMonth(birth.Month, birth.Year)
	if birth.Day > maxDay {
		sl.ReportError(birth.Day, "day", "Day", "validday", fmt.Sprintf("day must be between 1 and %d for month %d", maxDay, birth.Month))
	}
}

func getDaysInMonth(month, year int) int {
	// Get the last day of the month by going to first day of next month and subtracting one day
	lastDay := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC).Add(-24 * time.Hour)
	return lastDay.Day()
} 