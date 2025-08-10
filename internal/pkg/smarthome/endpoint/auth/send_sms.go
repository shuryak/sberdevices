package auth

import (
	"net/http"

	"github.com/shuryak/sberdevices/internal/pkg/pkce"
	"github.com/shuryak/sberdevices/internal/pkg/sbertypes"
	"github.com/shuryak/sberdevices/internal/pkg/smarthome/endpoint"
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
