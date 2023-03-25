package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"

	accessTokenExpiration = time.Hour
)

var (
	secret = []byte("+C3Nwz7LWIHMS7GdgdkV5tSOxzz++LocgdUdcjtSNpzMWfZ5JfoHElvWSOCs3MH9nPR4aWfTGvamRsY7icU4zQ==")
)

type Role string

func (u Role) IsSufficientToRole(role Role) bool {
	switch role {
	case RoleAdmin:
		if u == RoleAdmin {
			return true
		}
	case RoleUser:
		if u == RoleAdmin || u == RoleUser {
			return true
		}
	}
	return false
}

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	Role      Role
	CreatedAt time.Time
}

type Claims struct {
	UserID   uuid.UUID
	UserRole Role
}

type AccessToken struct {
	Claims     Claims
	SignedData string
	ExpiresAt  time.Time
}

func NewAccessToken(claims Claims) (*AccessToken, error) {
	expiration := time.Now().Add(accessTokenExpiration)
	signedData, err := signAccessToken(claims, expiration)
	if err != nil {
		return nil, fmt.Errorf("signing access token: %w", err)
	}

	return &AccessToken{
		SignedData: signedData,
		ExpiresAt:  expiration,
	}, nil
}

func (t AccessToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

type RefreshToken string

func NewRefreshToken() string {
	return uuid.New().String()
}

type TokenParser struct{}

func (t TokenParser) ParseAccessToken(data string) (*AccessToken, error) {
	tokenClaims, err := parseToken(data)
	if err != nil {
		return nil, err
	}

	return &AccessToken{
		Claims: Claims{
			UserID:   uuid.MustParse(tokenClaims.RegisteredClaims.Subject),
			UserRole: tokenClaims.CustomClaims.UserRole,
		},
		SignedData: data,
		ExpiresAt:  tokenClaims.ExpiresAt.Time,
	}, nil
}

type Session struct {
	AccessToken  AccessToken
	RefreshToken string
}

func NewSession(claims Claims) (*Session, error) {
	accessToken, err := NewAccessToken(claims)
	if err != nil {
		return nil, err
	}
	return &Session{
		AccessToken:  *accessToken,
		RefreshToken: NewRefreshToken(),
	}, nil
}

type customClaims struct {
	UserRole Role `json:"user_role"`
}

type claims struct {
	jwt.RegisteredClaims
	CustomClaims customClaims `json:"custom_claims"`
}

func signAccessToken(c Claims, expiresAt time.Time) (string, error) {
	jwtClaims := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   c.UserID.String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		CustomClaims: customClaims{
			UserRole: c.UserRole,
		},
	}

	signedData, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims).SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("new jwt with claims: %w", err)
	}

	return signedData, nil
}

func parseToken(data string) (*claims, error) {
	token, err := jwt.ParseWithClaims(data, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected jwt signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing jwt token with claims: %w", err)
	}

	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid jwt token")
	}

	return c, nil
}
