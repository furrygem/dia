package users

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repository Repository
}

func newService(pool *pgxpool.Pool) *service {
	return &service{
		repository: NewPostgresRepository(pool),
	}
}

func (s *service) checkUsernameAvailability(ctx context.Context, username string) (bool, error) {
	exists, err := s.repository.UsernameExists(username, ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *service) registerUser(ctx context.Context, userCreate UserCreateDTO) (*User, error) {
	usernameTaken, err := s.checkUsernameAvailability(ctx, userCreate.Username)
	if err != nil {
		return nil, fmt.Errorf("Could not verify username's availability")
	}

	if usernameTaken {
		return nil, fmt.Errorf("Username exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCreate.RawPassword), 10)
	if err != nil {
		return nil, err
	}
	user := User{}
	user.HashedPassword = fmt.Sprintf("%s", hashedPassword)
	return &user, nil
}

func (s *service) loginUser() {
}
