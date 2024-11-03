package devices

import (
	"errors"
	"net/http"

	"github.com/shuryak/sberdevices/internal/api"
	"github.com/shuryak/sberdevices/internal/transform"
	"github.com/shuryak/sberdevices/pkg/smarthome/client"
	"github.com/shuryak/sberdevices/pkg/yandex"
)

type queryReq struct {
	yandex.DevicesQueryRequest
}

func (req queryReq) Validate(_ *api.Context) error {
	return nil
}

func (h *Handlers) DevicesQuery(ctx *api.Context, req *queryReq) (*yandex.DevicesResponse, int) {
	ids := make([]string, len(req.Devices))
	for i := range req.Devices {
		ids[i] = req.Devices[i].ID
	}

	devices, err := h.client.GetDevices(getThirdPartyAccessToken(ctx), 5000, 0, ids...)
	if err != nil {
		if errors.Is(err, client.ErrTokenIsExpired) {
			return nil, http.StatusUnauthorized
		}

		return nil, http.StatusInternalServerError
	}

	return &yandex.DevicesResponse{
		RequestID: ctx.GetHeader("X-Request-Id"),
		Payload: &yandex.DevicesResponsePayload{
			Devices: transform.SberToYandexDevicesState(devices.Result),
		},
	}, http.StatusOK
}
