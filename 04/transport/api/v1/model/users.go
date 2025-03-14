package model

type User struct {
	Email string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Nickname  string `json:"nickname" validate:"omitempty,min=2"`
}
