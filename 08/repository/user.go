package repository

import (
	"context"

	"user-management-api/pkg/id"
	dbmodel "user-management-api/repository/sql/model"
	"user-management-api/repository/sql/query"
	"user-management-api/service/model"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) ReadUser(ctx context.Context, userID id.User) (*model.User, error) {
	var user dbmodel.User
	if err := pgxscan.Get(
		ctx,
		r.pool,
		&user,
		query.ReadUser,
		pgx.NamedArgs{
			"id": userID,
		},
	); err != nil {
		return nil, err
	}
	return &model.User{}, nil
}

func (r *UserRepository) ListUser(ctx context.Context) ([]model.User, error) {
	var users []dbmodel.User
	if err := pgxscan.Select(
		ctx,
		r.pool,
		&users,
		query.ListUser,
	); err != nil {
		return nil, err
	}
	response := make([]model.User, len(users))
	for i, user := range users {
		response[i] = model.User{
			ID: user.ID,
		}
	}
	return response, nil
}
