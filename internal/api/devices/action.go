package devices

import (
	"errors"
	"net/http"

	"github.com/shuryak/sberdevices/internal/api"
	"github.com/shuryak/sberdevices/internal/transform"
	"github.com/shuryak/sberdevices/pkg/smarthome/client"
	"github.com/shuryak/sberdevices/pkg/yandex"
)

type actionReq struct {
	yandex.DevicesActionRequest
}

func (req actionReq) Validate(_ *api.Context) error {
	return nil
}

func (h *Handlers) DevicesAction(ctx *api.Context, req *actionReq) (*yandex.DevicesResponse, int) {
	resp := &yandex.DevicesResponse{
		RequestID: ctx.GetHeader("X-Request-Id"),
		Payload:   &yandex.DevicesResponsePayload{},
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
		device := yandex.Device{
			ID: req.Payload.Devices[i].ID,
			ActionResult: &yandex.DeviceActionResult{
				Status: yandex.ResultStatusDone,
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

			device.Capabilities = append(device.Capabilities, yandex.DeviceCapability{
				Type: req.Payload.Devices[i].Capabilities[j].Type,
				State: &yandex.DeviceCapabilityState{
					Instance: req.Payload.Devices[i].Capabilities[j].State.Instance,
					ActionResult: &yandex.DeviceActionResult{
						Status: yandex.ResultStatusDone,
					},
				},
			})
		}

		resp.Payload.Devices = append(resp.Payload.Devices, device)
	}

	return resp, http.StatusOK
}
