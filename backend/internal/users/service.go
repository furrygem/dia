package users

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/furrygem/dia/internal/logging"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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

func (s *service) checkUsernameExists(ctx context.Context, username string) (bool, error) {
	exists, err := s.repository.UsernameExists(username, ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *service) registerUser(ctx context.Context, userCreate UserCreateDTO) (*User, error) {
	logger := logging.GetLogger()
	settings := getSettings()
	usernameTaken, err := s.checkUsernameExists(ctx, userCreate.Username)
	logger.Debugf("Username taken check for \"%s\": %t", userCreate.Username, usernameTaken)
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
	user.Username = userCreate.Username
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	user.Id = strings.ReplaceAll(id.String(), "-", "")
	logger.Debugf("Will set active: %t", settings.SetActiveAfterUserCreation)
	user.Active = settings.SetActiveAfterUserCreation
	user.CreatedAt = time.Now().Local()
	user.UpdatedAt = time.Now().Local()
	userInserted, err := s.repository.InsertUser(&user, ctx)
	if err != nil {
		return nil, err
	}
	logger.Infof("Inserted user %s with ID %s", userInserted.Username, userInserted.Id)
	return userInserted, nil
}

func (s *service) loginUser(u *User) jwt.Token {
	logger := logging.GetLogger()
	claims := jwt.RegisteredClaims{
		ExpiresAt: 10,
	}
}
