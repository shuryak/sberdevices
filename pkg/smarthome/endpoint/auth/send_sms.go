package auth

import (
	"net/http"

	"github.com/shuryak/sberdevices/pkg/pkce"
	"github.com/shuryak/sberdevices/pkg/sbertypes"
	"github.com/shuryak/sberdevices/pkg/smarthome/endpoint"
)

func SendSMS(phone string) (*endpoint.Endpoint, *pkce.Pair) {
	req, pkcePair := sbertypes.AuthDefaultSendSMSRequest(phone)

	return &endpoint.Endpoint{
		Method: "POST",
		Headers: http.Header{
			"Referer": {"SD"},
		},
		URL:  endpoint.MustURL("https://online.sberbank.ru/CSAFront/uapi/v2/authenticate"),
		Body: endpoint.MustJSON(req),
	}, pkcePair
}
