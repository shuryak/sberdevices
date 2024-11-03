package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shuryak/sberdevices/pkg/smarthome/endpoint"
	"github.com/shuryak/sberdevices/pkg/smarthome/util"
)

func (a *Authorizer) runEndpoint(ctx context.Context, e *endpoint.Endpoint, dest ...interface{}) error {
	switch statusCode, err := util.RunEndpoint(ctx, a.httpClient, e, dest...); {
	case err != nil:
		return err
	case statusCode != http.StatusOK:
		return fmt.Errorf("http status code %d", statusCode)
	}

	return nil
}
