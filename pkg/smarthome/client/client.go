package client

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/shuryak/sberdevices/pkg/sbertypes"
	"github.com/shuryak/sberdevices/pkg/smarthome/endpoint"
)

type Client struct {
	httpClient *http.Client
	timeout    time.Duration
	log        *log.Logger
}

func NewClient(
	timeout time.Duration,
	log *log.Logger,
) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		timeout:    timeout,
		log:        log,
	}
}

func (c *Client) SetDeviceState(accessToken, deviceID string, state ...*sbertypes.DeviceState) (*sbertypes.StateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp := &sbertypes.StateResponse{}
	err := c.runEndpoint(ctx, endpoint.State(accessToken, deviceID, state...), resp)
	return resp, err
}

func (c *Client) GetDevices(accessToken string, limit, offset int, ids ...string) (*sbertypes.DevicesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp := &sbertypes.DevicesResponse{}
	err := c.runEndpoint(ctx, endpoint.Devices(accessToken, limit, offset, ids...), resp)
	return resp, err
}

func (c *Client) GetDevice(accessToken, deviceID string) (*sbertypes.DeviceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp := &sbertypes.DeviceResponse{}
	err := c.runEndpoint(ctx, endpoint.Device(accessToken, deviceID), resp)
	return resp, err
}

var ErrTokenIsExpired = errors.New("token is expired")
