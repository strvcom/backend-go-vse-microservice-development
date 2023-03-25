package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type Service struct {
	users      map[uuid.UUID]*User
	usersMutex *sync.Mutex

	// Mapping between refresh token and user ID.
	refreshTokens      map[string]uuid.UUID
	refreshTokensMutex *sync.Mutex
}

func NewService() *Service {
	service := &Service{
		users:              make(map[uuid.UUID]*User),
		usersMutex:         &sync.Mutex{},
		refreshTokens:      make(map[string]uuid.UUID),
		refreshTokensMutex: &sync.Mutex{},
	}
	admin := &User{
		ID:        uuid.New(),
		Email:     "admin@gmail.com",
		Password:  "admin",
		Role:      RoleAdmin,
		CreatedAt: time.Now().UTC(),
	}
	service.users[admin.ID] = admin
	return service
}

func (s *Service) CreateSession(_ context.Context, user *User) (*Session, error) {
	session, err := NewSession(Claims{
		UserID:   user.ID,
		UserRole: user.Role,
	})
	if err != nil {
		return nil, err
	}

	s.refreshTokensMutex.Lock()
	defer s.refreshTokensMutex.Unlock()
	s.refreshTokens[session.RefreshToken] = user.ID

	return session, nil
}

func (s *Service) RefreshSession(_ context.Context, refreshToken string) (*Session, error) {
	s.refreshTokensMutex.Lock()
	defer s.refreshTokensMutex.Unlock()
	userID, ok := s.refreshTokens[refreshToken]
	if !ok {
		return nil, ErrInvalidRefreshToken
	}

	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()
	user := s.users[userID]

	session, err := NewSession(Claims{
		UserID:   user.ID,
		UserRole: user.Role,
	})
	if err != nil {
		return nil, err
	}

	delete(s.refreshTokens, refreshToken)
	s.refreshTokens[session.RefreshToken] = userID

	return session, nil
}

func (s *Service) DestroySession(_ context.Context, refreshToken string) {
	s.refreshTokensMutex.Lock()
	defer s.refreshTokensMutex.Unlock()
	delete(s.refreshTokens, refreshToken)
}

func (s *Service) CreateUser(_ context.Context, email, password string) (*User, error) {
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()

	if ok := s.existsUser(email); ok {
		return nil, ErrUserAlreadyExists
	}

	user := &User{
		ID:        uuid.New(),
		Email:     email,
		Password:  password,
		Role:      RoleUser,
		CreatedAt: time.Now().UTC(),
	}
	s.users[user.ID] = user

	return user, nil
}

func (s *Service) existsUser(email string) bool {
	for _, v := range s.users {
		if v.Email == email {
			return true
		}
	}
	return false
}

func (s *Service) ReadUser(_ context.Context, userID uuid.UUID) (*User, error) {
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()

	user, ok := s.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *Service) ReadUserByCredentials(_ context.Context, email, password string) (*User, error) {
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()

	for _, v := range s.users {
		if v.Email == email && v.Password == password {
			return v, nil
		}
	}
	return nil, ErrUserNotFound
}

func (s *Service) ReadUsers(_ context.Context) []User {
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()

	var users []User
	for _, u := range s.users {
		users = append(users, *u)
	}

	return users
}
