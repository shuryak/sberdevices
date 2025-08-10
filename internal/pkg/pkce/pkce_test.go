package pkce

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_generate(t *testing.T) {
	codeVerifier := []byte("fb0ffbaaf9b98e117b24e7d50ae083f7cd723dd6b1cf791cf82601b8")
	p := generate(codeVerifier)

	require.NotNil(t, p)
	require.Equal(t, p.CodeVerifier, string(codeVerifier))
	require.Equal(t, p.CodeChallenge, "Yp2tXQ_RzRTHClUd0a8OzDi5lqn_G-QC2X1B-3oYVvs")
}
