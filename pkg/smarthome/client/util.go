package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shuryak/sberhack/pkg/smarthome/endpoint"
	"github.com/shuryak/sberhack/pkg/smarthome/util"
)

func (c *Client) runEndpoint(ctx context.Context, e *endpoint.Endpoint, dest ...interface{}) error {
	switch statusCode, err := util.RunEndpoint(ctx, c.httpClient, e, dest...); {
	case err != nil:
		return err
	case statusCode == http.StatusUnauthorized:
		return ErrTokenIsExpired
	case statusCode != http.StatusOK:
		return fmt.Errorf("http status code %d", statusCode)
	}

	return nil
}
