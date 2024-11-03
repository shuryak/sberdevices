package pkce

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/shuryak/sberdevices/pkg/strrand"
)

type Pair struct {
	CodeVerifier  string
	CodeChallenge string
}

func GeneratePair() *Pair {
	return generate(strrand.RandSeq(64))
}

// https://base64.guru/standards/base64url
func generate(codeVerifier []byte) *Pair {
	codeVerifierChecksum := sha256.Sum256(codeVerifier)

	return &Pair{
		CodeVerifier: string(codeVerifier),
		CodeChallenge: strings.NewReplacer("+", "-", "/", "_", "=", "").
			Replace(base64.StdEncoding.EncodeToString(codeVerifierChecksum[:])),
	}
}
