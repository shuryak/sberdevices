package pkce

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/shuryak/sberdevices/internal/pkg/strrand"
)

type Pair struct {
	CodeVerifier  string
	CodeChallenge string
}

func GeneratePair() *Pair {
	return generate(strrand.RandSeq(64))
}

func generate(codeVerifier []byte) *Pair {
	codeVerifierChecksum := sha256.Sum256(codeVerifier)

	return &Pair{
		CodeVerifier:  string(codeVerifier),
		CodeChallenge: base64.RawURLEncoding.EncodeToString(codeVerifierChecksum[:]),
	}
}
