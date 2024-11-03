package model

import (
	"github.com/shuryak/sberdevices/pkg/pkce"
)

type AuthCodePayload struct {
	ThirdPartyPKCEPair        *pkce.Pair
	ThirdPartyOTPReceiverID   string
	ThirdPartyAuthOperationID string
}

func NewCodePayload(
	thirdPartyPKCEPair *pkce.Pair,
	thirdPartyOTPReceiverID, thirdPartyAuthOperationID string,
) *AuthCodePayload {
	return &AuthCodePayload{
		ThirdPartyPKCEPair:        thirdPartyPKCEPair,
		ThirdPartyOTPReceiverID:   thirdPartyOTPReceiverID,
		ThirdPartyAuthOperationID: thirdPartyAuthOperationID,
	}
}
