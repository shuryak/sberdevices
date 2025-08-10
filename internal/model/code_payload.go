package model

import (
	"github.com/shuryak/sberdevices/internal/pkg/pkce"
)

type AuthCodePayload struct {
	SmartHomePKCEPair        *pkce.Pair
	AccessToken              string
	SmartHomeOTPReceiverID   string
	SmartHomeAuthOperationID string
}
