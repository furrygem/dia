package pubkeys

import (
	"context"
	"fmt"

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

func (r *repo) GetByFingerprint(fingerprint []byte, ctx context.Context) (*PublicKey, error) {
	logger := logging.GetLogger()
	logger.Debugf("Retrieving key %X", fingerprint)
	logger.Info("executing")
	result := r.pool.QueryRow(ctx, "SELECT fingerprint, publickey FROM publickeys WHERE fingerprint = $1", fingerprint)
	pkey := &PublicKey{}
	var rawFingerprint []byte
	err := result.Scan(&rawFingerprint, &pkey.Key)
	pkey.Fingerprint = fmt.Sprintf("%X", rawFingerprint)
	if err != nil {
		logger.Debug(err)
		return nil, err
	}
	return pkey, nil
}

func (r *repo) StorePublicKey(fingerprint []byte, publickey string, ctx context.Context) (string, error) {
	_, err := r.pool.Exec(ctx, "INSERT INTO publickeys(fingerprint, publickey) VALUES ($1, $2)", fingerprint, publickey)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", fingerprint), nil
}
