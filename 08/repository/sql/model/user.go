package model

import (
	"time"

	"user-management-api/pkg/id"
)

type User struct {
	ID        id.User   `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
