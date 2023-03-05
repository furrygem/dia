package pubkeys

import "context"

type Reader interface {
	GetByFingerprint(fingerprint []byte, ctx context.Context) (*PublicKey, error)
}

type Writer interface {
	StorePublicKey(fingerpint []byte, publickey string, ctx context.Context) (string, error)
}

type Repository interface {
	Reader
	Writer
}
