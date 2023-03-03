package model

import "time"

type User struct {
	Email     string
	FirstName string
	LastName  string
	BirthDate BirthDate
}

type BirthDate struct {
	Day   int
	Month time.Month
	Year  int
}
