package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	*UserRepository
}

func New(pool *pgxpool.Pool) (*Repository, error) {
	return &Repository{
		UserRepository: NewUserRepository(pool),
	}, nil
}
