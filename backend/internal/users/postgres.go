package users

import (
	"context"

	"github.com/furrygem/dia/internal/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &repo{
		pool: pool,
	}
}

func (r *repo) GetUserByID(id string, ctx context.Context) (*User, error) {
	logger := logging.GetLogger()
	logger.Debug("Getting user by id %s", id)
	result := r.pool.QueryRow(ctx, "SELECT id, username, pfp_url, created_at, updated_at, active, hashed_password FROM users WHERE id = $1", id)
	user := User{}
	if err := result.Scan(&user.Id, &user.Username, &user.PfpURL, &user.CreatedAt, &user.UpdatedAt, &user.Active, &user.HashedPassword); err != nil {
		logger.Info(err.Error())
		return nil, err
	}
	return &user, nil
}
func (r *repo) GetUserByUsername(username string, ctx context.Context) (*User, error) {
	logger := logging.GetLogger()
	logger.Debug("Getting user by username %s", username)
	result := r.pool.QueryRow(ctx, "SELECT id, username, pfp_url, created_at, updated_at, active, hashed_password FROM users WHERE username = $1", username)
	user := User{}
	err := result.Scan(&user.Id, &user.Username, &user.PfpURL, &user.CreatedAt, &user.UpdatedAt, &user.Active, &user.HashedPassword)
	if err != nil {
		logger.Info(err)
		return nil, err
	}
	return &user, nil
}
func (r *repo) UsernameExists(username string, ctx context.Context) (bool, error) {
	logger := logging.GetLogger()
	logger.Debug("Checking if username %s exists", username)
	result := r.pool.QueryRow(ctx, "exists(SELECT 1 FROM users WHERE username = $1)", username)
	var exists bool
	err := result.Scan(&exists)
	if err != nil {
		logger.Info(err.Error())
		return false, err
	}
	return exists, nil
}

func (r *repo) InsertUser(username string, hashedPassword string, ctx context.Context) (*User, error) {
	logger := logging.GetLogger()
	logger.Debug("Inserting user with username=%s", username)
	result := r.pool.QueryRow(ctx, "INSERT INTO users (username, pfp_url, active, hashed_password) VALUES ($1, $2, $3, $4) RETURNING id, username, pfp_url, created_at, updated_at, active, hashed_pasword")
	user := &User{}
	err := result.Scan(&user.Id, &user.Username, &user.PfpURL, &user.CreatedAt, &user.UpdatedAt, &user.Active, &user.HashedPassword)
	if err != nil {
		logger.Info(err)
		return nil, err
	}
	return user, nil
}
func (r *repo) SetActive(id string, active bool, ctx context.Context) (bool, error) {
	logger := logging.GetLogger()
	logger.Debug("Setting active for user %s to %t", id, active)
	result := r.pool.QueryRow(ctx, "UPDATE users SET active = $1 WHERE id = $2 RETURNING active", active, id)
	err := result.Scan(&active)
	if err != nil {
		return active, err
	}
	return active, nil
}

func (r *repo) UpdatePassword(id string, hashedPassword string, ctx context.Context) (string, error) {
	logger := logging.GetLogger()
	logger.Debug("Updating password for user %s", id)
	result := r.pool.QueryRow(ctx, "UPDATE users SET hashed_password = $1 WHERE id = $2 RETURNING hashed_password", hashedPassword, id)
	err := result.Scan(&hashedPassword)
	if err != nil {
		return hashedPassword, err
	}
	return hashedPassword, nil
}
