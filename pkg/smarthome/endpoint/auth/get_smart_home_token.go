package auth

import (
	"fmt"
	"net/http"

	"github.com/shuryak/sberdevices/pkg/smarthome/endpoint"
)

func GetSmartHomeToken(csaFrontToken string) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		Method: "GET",
		URL:    endpoint.MustURL("https://mp-prom.salutehome.ru/v13/smarthome/token"),
		Headers: http.Header{
			"Host":          {"mp-prom.salutehome.ru"},
			"Authorization": {fmt.Sprintf("Bearer %s", csaFrontToken)},
		},
	}
}
