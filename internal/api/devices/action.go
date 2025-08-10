package devices

import (
	"errors"
	"net/http"

	"github.com/shuryak/sberdevices/internal/api"
	"github.com/shuryak/sberdevices/internal/pkg/smarthome/client"
	yandex2 "github.com/shuryak/sberdevices/internal/pkg/yandex"
	"github.com/shuryak/sberdevices/internal/transform"
)

type actionReq struct {
	yandex2.DevicesActionRequest
}

func (req actionReq) Validate(_ *api.Context) error {
	return nil
}

func (h *Handlers) DevicesAction(ctx *api.Context, req *actionReq) (*yandex2.DevicesResponse, int) {
	resp := &yandex2.DevicesResponse{
		RequestID: ctx.GetHeader("X-Request-Id"),
		Payload:   &yandex2.DevicesResponsePayload{},
	}

	ids := make([]string, len(req.Payload.Devices))
	for i := range req.Payload.Devices {
		ids[i] = req.Payload.Devices[i].ID
	}

	accessToken := getThirdPartyAccessToken(ctx)

	devices, err := h.client.GetDevices(accessToken, 5000, 0, ids...)
	if err != nil {
		if errors.Is(err, client.ErrTokenIsExpired) {
			return nil, http.StatusUnauthorized
		}

		return nil, http.StatusInternalServerError
	}
	devicesMap := devices.Result.ToMap()

	for i := range req.Payload.Devices {
		device := yandex2.Device{
			ID: req.Payload.Devices[i].ID,
			ActionResult: &yandex2.DeviceActionResult{
				Status: yandex2.ResultStatusDone,
			},
		}

		for j := range req.Payload.Devices[i].Capabilities {
			sberDeviceState := transform.YandexToSberDeviceState(
				devicesMap[req.Payload.Devices[i].ID].ReportedState.ToMap(),
				&req.Payload.Devices[i].Capabilities[j],
			)

			_, err := h.client.SetDeviceState(accessToken, req.Payload.Devices[i].ID, sberDeviceState...)
			if err != nil {
				if errors.Is(err, client.ErrTokenIsExpired) {
					return nil, http.StatusUnauthorized
				}

				return nil, http.StatusInternalServerError
			}

			device.Capabilities = append(device.Capabilities, yandex2.DeviceCapability{
				Type: req.Payload.Devices[i].Capabilities[j].Type,
				State: &yandex2.DeviceCapabilityState{
					Instance: req.Payload.Devices[i].Capabilities[j].State.Instance,
					ActionResult: &yandex2.DeviceActionResult{
						Status: yandex2.ResultStatusDone,
					},
				},
			})
		}

		resp.Payload.Devices = append(resp.Payload.Devices, device)
	}

	return resp, http.StatusOK
}
