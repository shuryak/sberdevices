package endpoint

import (
	"fmt"
	"net/http"
	"time"

	"github.com/shuryak/sberdevices/pkg/sbertypes"
)

func State(accessToken string, deviceID string, state ...*sbertypes.DeviceState) *Endpoint {
	return &Endpoint{
		Method: "PUT",
		URL:    MustURL(fmt.Sprintf("https://gateway.sbrdvc.xyz/gateway/v1/devices/%s/state", deviceID)),
		Headers: http.Header{
			"X-Auth-Jwt": {accessToken},
		},
		Body: MustJSON(
			sbertypes.StateRequest{
				Timestamp:    time.Now(),
				DeviceID:     deviceID,
				DesiredState: state,
			},
		),
	}
}
