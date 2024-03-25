package model

import (
	"time"
	"user-management-api/pkg/id"
)

type User struct {
	ID        id.User
	Email     string
	FirstName string
	LastName  string
	BirthDate BirthDate
}

type UpdateUserInput struct {
	LastName string
}

type BirthDate struct {
	Day   int
	Month time.Month
	Year  int
}
