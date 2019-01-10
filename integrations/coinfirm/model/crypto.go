package model

const (
	prefix = "address: `"
	suffix = "`"
)

// CryptoAddress represents the participant payment address.
type CryptoAddress string

// NewCryptoAddress constructs new CryptoAddress object.
func NewCryptoAddress(addr string) CryptoAddress {
	return CryptoAddress(prefix + addr + suffix)
}
