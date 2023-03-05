package pubkeys

type PublicKey struct {
	Fingerprint string `json:"fingerprint" db:"fingerprint"`
	Key         string `json:"key" db:"publickey"`
}
