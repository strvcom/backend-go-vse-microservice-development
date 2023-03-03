package model

import (
	"time"
	svcmodel "vse-course/service/model"
)

type User struct {
	Email     string    `json:"email" validate:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate BirthDate `json:"birth_date"`
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
