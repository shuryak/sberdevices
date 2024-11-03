package sbertypes

import (
	"github.com/shuryak/sberdevices/pkg/pkce"
	"github.com/shuryak/sberdevices/pkg/strrand"
)

const AuthSmartHomeClientID = "6835fd63-22c8-4c20-bd4a-2bba906afe5f"

type AuthSendSMSRequest struct {
	Authenticator AuthPayloadWrapper[AuthPayloadString]  `json:"authenticator"`
	Identifier    AuthPayloadWrapper[AuthPayloadString]  `json:"identifier"`
	Channel       AuthPayloadWrapper[AuthPayloadChannel] `json:"channel"`
}

func AuthDefaultSendSMSRequest(phone string) (*AuthSendSMSRequest, *pkce.Pair) {
	pkcePair := pkce.GeneratePair()

	return &AuthSendSMSRequest{
		Authenticator: AuthPayloadWrapper[AuthPayloadString]{
			Type: "sms_otp",
		},
		Identifier: AuthPayloadWrapper[AuthPayloadString]{
			Type: "phone",
			Data: &AuthPayloadString{
				Value: phone,
			},
		},
		Channel: AuthPayloadWrapper[AuthPayloadChannel]{
			Type: "web",
			Data: &AuthPayloadChannel{
				RSAData: AuthDefaultRSAData(),
				OIDC: &AuthOIDC{
					CodeChallengeMethod: "S256",
					Nonce:               strrand.RandSeqStr(64),
					Scope:               "openid",
					RedirectURI:         "homuzapp://host",
					CodeChallenge:       pkcePair.CodeChallenge,
					State:               strrand.RandSeqStr(64),
					ClientID:            AuthSmartHomeClientID,
					ResponseType:        "code",
				},
			},
		},
	}, pkcePair
}

type AuthSendSMSResponse struct {
	Authenticator []struct {
		Type     string `json:"type"`
		Lifetime int    `json:"lifetime"`
		Data     struct {
			Phones []string `json:"phones"`
		} `json:"data"`
		AttemptsRemaining      int  `json:"attempts_remaining"`
		InitializationRequired bool `json:"initialization_required"`
	} `json:"authenticator"`
	OUID string `json:"ouid"`
}

type AuthVerifySMSRequest struct {
	Authenticator AuthPayloadWrapper[AuthPayloadString]  `json:"authenticator"`
	Identifier    AuthPayloadWrapper[AuthPayloadString]  `json:"identifier"`
	Channel       AuthPayloadWrapper[AuthPayloadChannel] `json:"channel"`
}

func AuthDefaultVerifySMSRequest(ouid, otp string) *AuthVerifySMSRequest {
	return &AuthVerifySMSRequest{
		Authenticator: AuthPayloadWrapper[AuthPayloadString]{
			Type: "sms_otp",
			Data: &AuthPayloadString{
				Value: otp,
			},
		},
		Identifier: AuthPayloadWrapper[AuthPayloadString]{
			Type: "ouid",
			Data: &AuthPayloadString{
				Value: ouid,
			},
		},
		Channel: AuthPayloadWrapper[AuthPayloadChannel]{
			Type: "web",
			Data: &AuthPayloadChannel{
				RSAData: AuthDefaultRSAData(),
			},
		},
	}
}

type AuthVerifySMSResponse struct {
	ResponseData struct {
		RedirectUri string `json:"redirect_uri"`
		AuthCode    string `json:"authcode"`
		State       string `json:"state"`
	} `json:"response_data"`
}

type AuthGetCSAFrontTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	CodeVerifier string `json:"code_verifier"`
}

func AuthDefaultGetCSAFrontTokenRequest(authCode string, pkcePair *pkce.Pair) *AuthGetCSAFrontTokenRequest {
	return &AuthGetCSAFrontTokenRequest{
		GrantType:    "authorization_code",
		ClientID:     AuthSmartHomeClientID,
		Code:         authCode,
		RedirectURI:  "homuzapp://host",
		CodeVerifier: pkcePair.CodeVerifier,
	}
}

type AuthGetCSAFrontTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    uint64 `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthRefreshCSAFrontTokenRequest struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

func AuthDefaultRefreshCSAFrontTokenRequest(refreshToken string) *AuthRefreshCSAFrontTokenRequest {
	return &AuthRefreshCSAFrontTokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}
}

type AuthRefreshCSAFrontTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    uint64 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type AuthGetSmartHomeTokenResponse struct {
	State struct {
		Status int `json:"status"`
	} `json:"state"`
	Token string `json:"token"`
}

type AuthPayloadWrapper[T interface{}] struct {
	Type string `json:"type"`
	Data *T     `json:"data,omitempty"`
}

type AuthPayloadString struct {
	Value string `json:"value"`
}

type AuthPayloadChannel struct {
	RSAData *AuthRSAData `json:"rsa_data,omitempty"`
	OIDC    *AuthOIDC    `json:"oidc,omitempty"`
}

type AuthRSAData struct {
	DevicePrint           string `json:"deviceprint"`
	HTMLInjection         string `json:"htmlinjection"`
	DOMElements           string `json:"dom_elements"`
	ManVSMachineDetection string `json:"manvsmachinedetection"`
	JSEvents              string `json:"js_events"`
}

func AuthDefaultRSAData() *AuthRSAData {
	return &AuthRSAData{
		DevicePrint:           "{}",
		HTMLInjection:         "htmlinjection",
		DOMElements:           "dom_elements",
		ManVSMachineDetection: "manvsmachinedetection",
		JSEvents:              "js_events",
	}
}

type AuthOIDC struct {
	CodeChallengeMethod string `json:"code_challenge_method"`
	Nonce               string `json:"nonce"`
	Scope               string `json:"scope"`
	RedirectURI         string `json:"redirect_uri"`
	CodeChallenge       string `json:"code_challenge"`
	State               string `json:"state"`
	ClientID            string `json:"client_id"`
	ResponseType        string `json:"response_type"`
}
