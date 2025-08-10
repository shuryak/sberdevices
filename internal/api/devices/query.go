package devices

import (
	"errors"
	"net/http"

	"github.com/shuryak/sberdevices/internal/api"
	"github.com/shuryak/sberdevices/internal/pkg/smarthome/client"
	yandex2 "github.com/shuryak/sberdevices/internal/pkg/yandex"
	"github.com/shuryak/sberdevices/internal/transform"
)

type queryReq struct {
	yandex2.DevicesQueryRequest
}

func (req queryReq) Validate(_ *api.Context) error {
	return nil
}

func (h *Handlers) DevicesQuery(ctx *api.Context, req *queryReq) (*yandex2.DevicesResponse, int) {
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

	return &yandex2.DevicesResponse{
		RequestID: ctx.GetHeader("X-Request-Id"),
		Payload: &yandex2.DevicesResponsePayload{
			Devices: transform.SberToYandexDevicesState(devices.Result),
		},
	}, http.StatusOK
}
