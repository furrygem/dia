package pubkeys

import (
	"fmt"
)

type PublicKey struct {
	rawFingerprint string
	Fingerprint    string `json:"fingerprint" db:"fingerprint"`
	Key            string `json:"key" db:"publickey"`
}

func (pk *PublicKey) Prepare() {
	pk.Fingerprint = fmt.Sprintf("%X", pk.rawFingerprint)
}
