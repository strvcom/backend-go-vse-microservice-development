package v1

import (
	"context"

	"user-management-api/pkg/id"
	svcmodel "user-management-api/service/model"
)

type Service interface {
	CreateUser(ctx context.Context, user svcmodel.User) error
	ListUsers(ctx context.Context) ([]svcmodel.User, error)
	GetUser(ctx context.Context, userID id.User) (*svcmodel.User, error)
	UpdateUser(ctx context.Context, userID id.User, updateUserInput svcmodel.UpdateUserInput) (*svcmodel.User, error)
	DeleteUser(ctx context.Context, userID id.User) error
}
