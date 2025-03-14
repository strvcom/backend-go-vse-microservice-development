package model

type User struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Nickname  string `json:"nickname" validate:"omitempty,min=2"`
	BirthDate BirthDate `json:"birth_date" validate:"required"`
}

type BirthDate struct {
	Day   int        `json:"day" validate:"required,min=1,max=31"`
	Month int        `json:"month" validate:"required,min=1,max=12"`
	Year  int        `json:"year" validate:"required,min=1900,max=2100"`
}
