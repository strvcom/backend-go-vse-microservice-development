package service

import (
	"context"

	"vse-course/service/errors"
	"vse-course/service/model"
)

var (
	users = map[string]model.User{}
)

// CreateUser saves user in map under email as a key.
func (Service) CreateUser(_ context.Context, user model.User) error {
	if _, exists := users[user.Email]; exists {
		return errors.ErrUserAlreadyExists
	}

	users[user.Email] = user

	return nil
}

// ListUsers returns list of users in array of users.
func (Service) ListUsers(_ context.Context) []model.User {
	usersList := make([]model.User, 0, len(users))
	for _, user := range users {
		usersList = append(usersList, user)
	}

	return usersList
}

// GetUser returns an user with specified email.
func (Service) GetUser(_ context.Context, email string) (model.User, error) {
	user, exists := users[email]

	if !exists {
		return model.User{}, errors.ErrUserDoesntExists
	}

	return user, nil
}

// UpdateUser updates attributes of a specified user.
func (Service) UpdateUser(_ context.Context, email string, user model.User) (model.User, error) {
	oldUser, exists := users[email]

	if !exists {
		return model.User{}, errors.ErrUserDoesntExists
	}

	if oldUser.Email == user.Email {
		users[email] = user
	} else {
		users[user.Email] = user

		delete(users, email)
	}

	return user, nil
}

// DeleteUser deletes user from memory.
func (Service) DeleteUser(_ context.Context, email string) error {
	if _, exists := users[email]; !exists {
		return errors.ErrUserDoesntExists
	}

	delete(users, email)

	return nil
}
