package endpoint

import (
	"fmt"
	"net/http"
)

func Device(accessToken, deviceID string) *Endpoint {
	return &Endpoint{
		Method: "GET",
		Headers: http.Header{
			"X-Auth-Jwt": {accessToken},
		},
		URL: MustURL(fmt.Sprintf("https://gateway.sbrdvc.xyz/gateway/v1/devices/%s", deviceID)),
	}
}
