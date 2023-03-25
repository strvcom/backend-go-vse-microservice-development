package util

import (
	"context"

	"auth/custom/service"

	"github.com/google/uuid"
)

var (
	contextKey = struct {
		userID   ctxKeyUserID
		userRole ctxKeyUserRole
	}{}
)

type (
	ctxKeyUserID   struct{}
	ctxKeyUserRole struct{}
)

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, contextKey.userID, userID)
}

func UserIDFromCtx(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(contextKey.userID).(uuid.UUID)
	return userID, ok
}

func WithUserRole(ctx context.Context, role service.Role) context.Context {
	return context.WithValue(ctx, contextKey.userRole, role)
}

func UserRoleFromCtx(ctx context.Context) (service.Role, bool) {
	userRole, ok := ctx.Value(contextKey.userRole).(service.Role)
	return userRole, ok
}
