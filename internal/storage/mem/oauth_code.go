package mem

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/shuryak/sberhack/internal/model"
	"github.com/shuryak/sberhack/pkg/pkce"
	"github.com/shuryak/sberhack/pkg/strrand"
)

type OAuthCodeStorage struct {
	codeLength int
	storage    sync.Map
}

func NewOAuthCodeStorage(codeLength int) *OAuthCodeStorage {
	return &OAuthCodeStorage{codeLength: codeLength}
}

func (s *OAuthCodeStorage) Create(
	_ context.Context,
	thirdPartyPKCEPair *pkce.Pair,
	thirdPartyOTPReceiverID, thirdPartyAuthOperationID string,
) (authCode string, err error) {
	code := strrand.RandSeqStr(s.codeLength)
	if len(code) != s.codeLength {
		return "", fmt.Errorf("invalid random string generation with length %d", s.codeLength)
	}

	s.storage.Store(code, *model.NewCodePayload(thirdPartyPKCEPair, thirdPartyOTPReceiverID, thirdPartyAuthOperationID))

	return code, err
}

func (s *OAuthCodeStorage) Get(_ context.Context, code string) (*model.AuthCodePayload, error) {
	v, ok := s.storage.Load(code)
	if !ok {
		return nil, ErrCodeNotFound
	}

	codePayload := v.(model.AuthCodePayload)

	return &codePayload, nil
}

func (s *OAuthCodeStorage) Delete(_ context.Context, code string) error {
	_, ok := s.storage.Load(code)
	if !ok {
		return ErrCodeNotFound
	}

	s.storage.Delete(code)

	return nil
}

var (
	ErrCodeNotFound = errors.New("code not found")
)
