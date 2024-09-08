package api

import (
	"net/http"
)

func CORS(ctx *Context) {
	ctx.SetHeader("Access-Control-Allow-Origin", "*")
	ctx.SetHeader("Access-Control-Allow-ThirdPartyCredentials", "true")
	ctx.SetHeader(
		"Access-Control-Allow-Headers",
		"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
	)
	ctx.SetHeader("Access-Control-Allow-Methods", "POST, OPTIONS, GET")

	if ctx.GetMethod() == "OPTIONS" {
		_ = ctx.WriteResponse(http.StatusOK, nil)
		ctx.StopChain()
	}
}
