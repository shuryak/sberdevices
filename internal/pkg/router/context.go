package router

import (
	"context"
	"net/http"
)

type Context interface {
	context.Context
	SetCancellableCtx(ctx context.Context, cancel context.CancelFunc)
	SetHTTPWriter(w http.ResponseWriter)
	SetHTTPRequest(r *http.Request)
	SetErrHandler(errHandler func(err error) (interface{}, int))
	Decode(dest interface{}) error
	WriteResponse(resp interface{}, statusCode int) error
}
