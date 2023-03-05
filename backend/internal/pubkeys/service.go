package pubkeys

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/furrygem/dia/internal/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type service struct {
	repository Repository
}

func newService(pool *pgxpool.Pool) *service {
	return &service{
		repository: NewPostgresRepository(pool),
	}
}

func (s *service) writePublicKey(requestBody []byte, ctx context.Context) (string, error) {
	logger := logging.GetLogger()
	logger.Debugf("Reading body %x into keyring", requestBody)
	keyring, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(requestBody))
	if err != nil {
		return "", err
	}
	var fingerprint []byte
	for _, key := range keyring {
		logger.Debugf("Fingerprint: %x", fingerprint)
		fingerprint = key.PrimaryKey.Fingerprint
	}

	fprinthex := fmt.Sprintf("%X", fingerprint)
	_, err = s.repository.StorePublicKey(fingerprint, fmt.Sprintf("%s", requestBody), ctx)
	if err != nil {
		return fprinthex, err
	}
	return fprinthex, nil
}
