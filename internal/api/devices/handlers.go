package devices

import (
	"github.com/shuryak/sberdevices/internal/api"
	"github.com/shuryak/sberdevices/internal/api/common"
	"github.com/shuryak/sberdevices/internal/oauth"
	"github.com/shuryak/sberdevices/pkg/smarthome/client"
)

type Handlers struct {
	client *client.Client
	flow   *oauth.CodeFlowWithOTP
}

func NewHandlers(client *client.Client, flow *oauth.CodeFlowWithOTP) *Handlers {
	return &Handlers{
		client: client,
		flow:   flow,
	}
}

func getAccessToken(ctx *api.Context) string {
	return ctx.Value(common.CtxKeyAccessToken).(string)
}

func getThirdPartyAccessToken(ctx *api.Context) string {
	return ctx.Value(common.CtxKeyThirdPartyAccessToken).(string)
}

func getThirdPartyRefreshToken(ctx *api.Context) string {
	return ctx.Value(common.CtxKeyThirdPartyRefreshToken).(string)
}
