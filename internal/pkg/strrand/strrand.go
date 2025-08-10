package strrand

import (
	"crypto/rand"
	"math/big"
)

func RandSeq(length uint) []byte {
	const charset = "YQePtiF8UvdIupcgS9yWkZ7E1HnLBxOJNjRbKo5VlMAwfhG6mDTs2XC0z43raq"

	seq := make([]byte, length)

	for i := uint(0); i < length; i++ {
		randIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return []byte{}
		}

		seq[i] = charset[randIdx.Int64()]
	}

	return seq
}

func RandSeqStr(length uint) string {
	return string(RandSeq(length))
}
