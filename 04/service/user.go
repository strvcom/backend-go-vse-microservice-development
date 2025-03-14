package service

import (
	"context"

	"user-management-api/service/errors"
	"user-management-api/service/model"
	"user-management-api/transport/util"
)

var (
	users = map[string]model.User{}
)

// CreateUser saves user in map under email as a key.
func (Service) CreateUser(ctx context.Context, user model.User) error {
	logger := util.NewServerLogger("service")
	logger.Info("creating user")

	if _, exists := users[user.Email]; exists {
		logger.Error("user already exists", errors.ErrUserAlreadyExists)
		return errors.ErrUserAlreadyExists
	}

	users[user.Email] = user
	logger.Info("user created successfully")
	return nil
}

// ListUsers returns list of users in array of users.
func (Service) ListUsers(ctx context.Context) []model.User {
	logger := util.NewServerLogger("service")
	logger.Info("listing users")

	usersList := make([]model.User, 0, len(users))
	for _, user := range users {
		usersList = append(usersList, user)
	}

	logger.Info("users retrieved successfully")
	return usersList
}

// GetUser returns an user with specified email.
func (Service) GetUser(ctx context.Context, email string) (model.User, error) {
	logger := util.NewServerLogger("service")
	logger.Info("fetching user")

	user, exists := users[email]
	if !exists {
		logger.Error("user not found", errors.ErrUserDoesntExists)
		return model.User{}, errors.ErrUserDoesntExists
	}

	logger.Info("user retrieved successfully")
	return user, nil
}

// UpdateUser updates attributes of a specified user.
func (Service) UpdateUser(ctx context.Context, email string, user model.User) (model.User, error) {
	logger := util.NewServerLogger("service")
	logger.Info("updating user")

	oldUser, exists := users[email]
	if !exists {
		logger.Error("user not found", errors.ErrUserDoesntExists)
		return model.User{}, errors.ErrUserDoesntExists
	}

	if oldUser.Email == user.Email {
		users[email] = user
	} else {
		users[user.Email] = user
		delete(users, email)
	}

	logger.Info("user updated successfully")
	return user, nil
}

// DeleteUser deletes user from memory.
func (Service) DeleteUser(ctx context.Context, email string) error {
	logger := util.NewServerLogger("service")
	logger.Info("deleting user")

	if _, exists := users[email]; !exists {
		logger.Error("user not found", errors.ErrUserDoesntExists)
		return errors.ErrUserDoesntExists
	}

	delete(users, email)
	logger.Info("user deleted successfully")
	return nil
}
