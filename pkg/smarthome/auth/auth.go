package auth

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/shuryak/sberdevices/pkg/pkce"
	"github.com/shuryak/sberdevices/pkg/sbertypes"
	"github.com/shuryak/sberdevices/pkg/smarthome/endpoint/auth"
)

type Authorizer struct {
	httpClient *http.Client
	log        *log.Logger
}

func NewAuthorizer(log *log.Logger) *Authorizer {
	// TODO: only https://online.sberbank.ru:4431/CSAFront/api/service/oidc/v3/token
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Authorizer{
		httpClient: &http.Client{Transport: tr},
		log:        log,
	}
}

func (a *Authorizer) SendOTP(ctx context.Context, phone string) (pkcePair *pkce.Pair, ouid string, err error) {
	a.log.Println("send otp method called")

	req, pkcePair := auth.SendSMS(phone)

	resp := &sbertypes.AuthSendSMSResponse{}
	err = a.runEndpoint(ctx, req, resp)
	if err != nil {
		return nil, "", err
	}

	return pkcePair, resp.OUID, nil
}

func (a *Authorizer) GetSmartHomeTokenByOTP(
	ctx context.Context, ouid string, pkcePair *pkce.Pair, otp string,
) (*TokenResult, error) {
	a.log.Println("get smart home token by otp method called")

	verifySMSResp := &sbertypes.AuthVerifySMSResponse{}
	err := a.runEndpoint(ctx, auth.VerifySMS(ouid, otp), verifySMSResp)
	if err != nil {
		return nil, err
	}

	getCSAFrontTokenResp := &sbertypes.AuthGetCSAFrontTokenResponse{}
	err = a.runEndpoint(ctx, auth.GetCSAFrontToken(verifySMSResp.ResponseData.AuthCode, pkcePair),
		getCSAFrontTokenResp,
	)
	if err != nil {
		return nil, err
	}

	getSmartHomeTokenResp := &sbertypes.AuthGetSmartHomeTokenResponse{}
	err = a.runEndpoint(ctx, auth.GetSmartHomeToken(getCSAFrontTokenResp.AccessToken), getSmartHomeTokenResp)
	if err != nil {
		return nil, err
	}

	return &TokenResult{
		CSAFrontRefreshToken:   getCSAFrontTokenResp.RefreshToken,
		CSAFrontAccessTokenTTL: time.Second * time.Duration(getCSAFrontTokenResp.ExpiresIn),
		SmartHomeToken:         getSmartHomeTokenResp.Token,
	}, nil
}

func (a *Authorizer) RefreshSmartHomeToken(ctx context.Context, csaFrontRefreshToken string) (*TokenResult, error) {
	a.log.Println("refresh smart home token method called")

	refreshCSAFrontTokenResp := &sbertypes.AuthRefreshCSAFrontTokenResponse{}
	err := a.runEndpoint(ctx, auth.RefreshCSAFrontToken(csaFrontRefreshToken), refreshCSAFrontTokenResp)
	if err != nil {
		return nil, err
	}

	getSmartHomeTokenResp := &sbertypes.AuthGetSmartHomeTokenResponse{}
	err = a.runEndpoint(ctx, auth.GetSmartHomeToken(refreshCSAFrontTokenResp.AccessToken), getSmartHomeTokenResp)
	if err != nil {
		return nil, err
	}

	return &TokenResult{
		CSAFrontRefreshToken:   refreshCSAFrontTokenResp.RefreshToken,
		CSAFrontAccessTokenTTL: time.Second * time.Duration(refreshCSAFrontTokenResp.ExpiresIn),
		SmartHomeToken:         getSmartHomeTokenResp.Token,
	}, nil
}

type TokenResult struct {
	CSAFrontRefreshToken   string
	CSAFrontAccessTokenTTL time.Duration
	SmartHomeToken         string
}
