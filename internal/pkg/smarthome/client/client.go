package client

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	sbertypes2 "github.com/shuryak/sberdevices/internal/pkg/sbertypes"
	endpoint2 "github.com/shuryak/sberdevices/internal/pkg/smarthome/endpoint"
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

func (c *Client) SetDeviceState(accessToken, deviceID string, state ...*sbertypes2.DeviceState) (*sbertypes2.StateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp := &sbertypes2.StateResponse{}
	err := c.runEndpoint(ctx, endpoint2.State(accessToken, deviceID, state...), resp)
	return resp, err
}

func (c *Client) GetDevices(accessToken string, limit, offset int, ids ...string) (*sbertypes2.DevicesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp := &sbertypes2.DevicesResponse{}
	err := c.runEndpoint(ctx, endpoint2.Devices(accessToken, limit, offset, ids...), resp)
	return resp, err
}

func (c *Client) GetDevice(accessToken, deviceID string) (*sbertypes2.DeviceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp := &sbertypes2.DeviceResponse{}
	err := c.runEndpoint(ctx, endpoint2.Device(accessToken, deviceID), resp)
	return resp, err
}

var ErrTokenIsExpired = errors.New("token is expired")
