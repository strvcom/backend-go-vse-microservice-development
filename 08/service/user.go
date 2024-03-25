package service

import (
	"context"

	"user-management-api/pkg/id"
	"user-management-api/service/model"
)

// CreateUser saves user in map under email as a key.
func (Service) CreateUser(_ context.Context, user model.User) error {
	panic("not implemented")
}

// ListUsers returns list of users in array of users.
func (s Service) ListUsers(ctx context.Context) ([]model.User, error) {
	users, err := s.repository.ListUser(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUser returns an user with specified email.
func (s Service) GetUser(ctx context.Context, userID id.User) (*model.User, error) {
	user, err := s.repository.ReadUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates attributes of a specified user.
func (Service) UpdateUser(_ context.Context, userID id.User, updateUserInput model.UpdateUserInput) (*model.User, error) {
	panic("not implemented")
}

// DeleteUser deletes user from memory.
func (Service) DeleteUser(_ context.Context, userID id.User) error {
	panic("not implemented")
}
