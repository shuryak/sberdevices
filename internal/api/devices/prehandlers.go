package devices

import (
	"context"
	"strings"

	"github.com/shuryak/sberhack/internal/api"
	"github.com/shuryak/sberhack/internal/api/common"
)

func (h *Handlers) Tokens(ctx *api.Context) {
	accessToken := ctx.GetHeader("Authorization")
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	session, err := h.flow.GetSessionByAccessToken(ctx, accessToken)
	if err != nil {
		// TODO: handle
		return
	}

	// TODO: good functions to get this ctx values
	ctx.Context = context.WithValue(ctx.Context, common.CtxKeyAccessToken, session.AccessToken)
	ctx.Context = context.WithValue(ctx.Context, common.CtxKeyThirdPartyAccessToken, session.ThirdPartyAccessToken)
	ctx.Context = context.WithValue(ctx.Context, common.CtxKeyThirdPartyRefreshToken, session.ThirdPartyRefreshToken)
}
