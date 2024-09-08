package devices

import (
	"errors"
	"net/http"

	"github.com/shuryak/sberhack/internal/api"
	"github.com/shuryak/sberhack/internal/transform"
	"github.com/shuryak/sberhack/pkg/smarthome/client"
	"github.com/shuryak/sberhack/pkg/yandex"
)

type listReq struct{}

func (req listReq) Validate(_ *api.Context) error {
	return nil
}

func (h *Handlers) Devices(ctx *api.Context, _ *listReq) (*yandex.DevicesResponse, int) {
	devices, err := h.client.GetDevices(getThirdPartyAccessToken(ctx), 5000, 0)
	if err != nil {
		if errors.Is(err, client.ErrTokenIsExpired) {
			return nil, http.StatusUnauthorized
		}

		return nil, http.StatusInternalServerError
	}

	return &yandex.DevicesResponse{
		RequestID: ctx.GetHeader("X-Request-Id"),
		Payload: &yandex.DevicesResponsePayload{
			UserID:  "custom_user_id", // TODO: user_id
			Devices: transform.SberToYandexDevices(devices.Result),
		},
	}, http.StatusOK
}
