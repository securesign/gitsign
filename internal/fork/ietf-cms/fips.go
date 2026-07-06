package cms

import (
	"crypto"
	"crypto/fips140"
	"fmt"
)

// RHTAS FIPS - DO NOT REMOVE
// ========================================

func checkFIPSHash(hash crypto.Hash) error {
	if fips140.Enabled() && !isFIPSApprovedHash(hash) {
		return fmt.Errorf("hash algorithm %v is not approved in FIPS 140-3 mode", hash)
	}
	return nil
}

func isFIPSApprovedHash(h crypto.Hash) bool {
	switch h {
	case crypto.SHA224, crypto.SHA256, crypto.SHA384, crypto.SHA512,
		crypto.SHA3_224, crypto.SHA3_256, crypto.SHA3_384, crypto.SHA3_512:
		return true
	default:
		return false
	}
}

// ========================================
