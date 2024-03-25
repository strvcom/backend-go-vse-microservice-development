package service

import (
	"context"

	"user-management-api/pkg/id"
	"user-management-api/service/model"
)

type Repository interface {
	ReadUser(ctx context.Context, userID id.User) (*model.User, error)
	ListUser(ctx context.Context) ([]model.User, error)
}

type Service struct {
	repository Repository
}

func NewService(
	repository Repository,
) (Service, error) {
	return Service{
		repository: repository,
	}, nil
}
