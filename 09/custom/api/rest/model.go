package rest

import (
	"time"

	"auth/custom/service"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

func FromUser(user *service.User) User {
	return User{
		ID:        user.ID,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}
}

type CreateUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateUserResp struct {
	User    User    `json:"user"`
	Session Session `json:"session"`
}

type CreateSessionInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Session struct {
	AccessToken          string    `json:"accessToken"`
	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt"`
	RefreshToken         string    `json:"refreshToken"`
}

func FromSession(session *service.Session) Session {
	return Session{
		AccessToken:          session.AccessToken.SignedData,
		AccessTokenExpiresAt: session.AccessToken.ExpiresAt,
		RefreshToken:         session.RefreshToken,
	}
}

type RefreshSessionInput struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
