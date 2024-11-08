package auth

import (
	"net/http"

	"github.com/shuryak/sberdevices/pkg/sbertypes"
	"github.com/shuryak/sberdevices/pkg/smarthome/endpoint"
)

func VerifySMS(ouid, otp string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Method: "POST",
		Headers: http.Header{
			"Referer": {"SD"},
		},
		URL:  endpoint.MustURL("https://online.sberbank.ru/CSAFront/uapi/v2/verify"),
		Body: endpoint.MustJSON(sbertypes.AuthDefaultVerifySMSRequest(ouid, otp)),
	}
}
