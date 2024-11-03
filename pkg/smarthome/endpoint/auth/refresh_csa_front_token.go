package auth

import (
	"net/http"
	"net/url"

	"github.com/shuryak/sberdevices/pkg/sbertypes"
	"github.com/shuryak/sberdevices/pkg/smarthome/endpoint"
)

func RefreshCSAFrontToken(refreshToken string) *endpoint.Endpoint {
	req := sbertypes.AuthDefaultRefreshCSAFrontTokenRequest(refreshToken)

	return &endpoint.Endpoint{
		Method: "POST",
		URL:    endpoint.MustURL("https://online.sberbank.ru:4431/CSAFront/api/service/oidc/v3/token"),
		Headers: http.Header{
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
		Body: []byte(url.Values{
			"grant_type":    []string{req.GrantType},
			"refresh_token": []string{req.RefreshToken},
		}.Encode()),
	}
}
