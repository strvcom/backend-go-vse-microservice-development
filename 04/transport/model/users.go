package model

import (
	"errors"
	"time"

	svcmodel "vse-course/service/model"
)

type User struct {
	Email     string    `json:"email" validate:"email"`
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name"`
	BirthDate BirthDate `json:"birth_date"`
}

func (bd BirthDate) ValidateBirthDate() error {
	if bd.Month <= 0 || bd.Month > 12 {
		return errors.New("invalid month")
	}

	if bd.Day < 0 {
		return errors.New("invalid day")
	}

	switch bd.Month {
	case 1, 3, 5, 7, 8, 10, 12:
		if bd.Day > 31 {
			return errors.New("invalid day")
		}
	case 4, 6, 9, 11:
		if bd.Day > 30 {
			return errors.New("invalid day")
		}
	default:
		if bd.Day > 28 {
			return errors.New("invalid day")
		}
	}

	return nil
}

func ToSvcUser(u User) svcmodel.User {
	return svcmodel.User{
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		BirthDate: svcmodel.BirthDate{
			Day:   u.BirthDate.Day,
			Month: time.Month(u.BirthDate.Month),
			Year:  u.BirthDate.Year,
		},
	}
}

func ToNetUser(u svcmodel.User) User {
	return User{
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		BirthDate: BirthDate{
			Day:   u.BirthDate.Day,
			Month: int(u.BirthDate.Month),
			Year:  u.BirthDate.Year,
		},
	}
}

func ToNetUsers(users []svcmodel.User) []User {
	netUsers := make([]User, 0, len(users))
	for _, user := range users {
		netUsers = append(netUsers, ToNetUser(user))
	}
	return netUsers
}

type BirthDate struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}
