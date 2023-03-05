package users

import "github.com/jackc/pgx/v5/pgxpool"

type service struct {
	repository Repository
}

func newService(pool *pgxpool.Pool) *service {
	return &service{
		repository: NewPostgresRepository(pool),
	}
}

func (s *service) registerUser() {}

func (s *service) loginUser() {}
