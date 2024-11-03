package auth

import (
	"net/http"
	"net/url"

	"github.com/shuryak/sberdevices/pkg/pkce"
	"github.com/shuryak/sberdevices/pkg/sbertypes"
	"github.com/shuryak/sberdevices/pkg/smarthome/endpoint"
)

func GetCSAFrontToken(authCode string, pkcePair *pkce.Pair) *endpoint.Endpoint {
	req := sbertypes.AuthDefaultGetCSAFrontTokenRequest(authCode, pkcePair)

	return &endpoint.Endpoint{
		Method: "POST",
		URL:    endpoint.MustURL("https://online.sberbank.ru:4431/CSAFront/api/service/oidc/v3/token"),
		Headers: http.Header{
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
		Body: []byte(url.Values{
			"grant_type":    []string{req.GrantType},
			"client_id":     []string{req.ClientID},
			"code":          []string{req.Code},
			"redirect_uri":  []string{req.RedirectURI},
			"code_verifier": []string{req.CodeVerifier},
		}.Encode()),
	}
}
