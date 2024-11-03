package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/shuryak/sberdevices/pkg/query"
	"github.com/shuryak/sberdevices/pkg/router"
)

type Context struct {
	context.Context
	cancel context.CancelFunc
	w      http.ResponseWriter
	r      *http.Request
}

// Check for implementation
var _ router.Context = (*Context)(nil)

type Validator interface {
	Validate(ctx *Context) error
}

func (ctx *Context) SetCancellableCtx(baseCtx context.Context, cancel context.CancelFunc) {
	ctx.Context = baseCtx
	ctx.cancel = cancel
}

func (ctx *Context) SetHTTPWriter(w http.ResponseWriter) {
	ctx.w = w
}

func (ctx *Context) SetHTTPRequest(r *http.Request) {
	ctx.r = r
}

func (ctx *Context) StopChain() {
	ctx.cancel()
}

func (ctx *Context) GetMethod() string {
	return ctx.r.Method
}

func (ctx *Context) BodyBytes() ([]byte, error) {
	return io.ReadAll(ctx.r.Body)
}

func (ctx *Context) Query() url.Values {
	return ctx.r.URL.Query()
}

func (ctx *Context) Redirect(url string) {
	ctx.StopChain()
	http.Redirect(ctx.w, ctx.r, url, http.StatusFound)
}

func (ctx *Context) Decode(dest interface{}) error {
	var decoder interface {
		Decode(interface{}) error
	}

	switch ctx.GetHeader("Content-Type") {
	case "application/x-www-form-urlencoded":
		bytes, err := ctx.BodyBytes()
		if err != nil {
			return err
		}

		values, err := url.ParseQuery(string(bytes))
		if err != nil {
			return err
		}

		decoder = query.NewDecoder(values)
	case "application/json":
		decoder = json.NewDecoder(ctx.r.Body)
	default:
		decoder = query.NewDecoder(ctx.r.URL.Query())
	}

	err := decoder.Decode(dest)
	if err != nil {
		return err
	}

	return dest.(Validator).Validate(ctx)
}

func (ctx *Context) SetHeader(key string, value string) {
	ctx.w.Header().Set(key, value)
}

func (ctx *Context) GetHeader(key string) string {
	return ctx.r.Header.Get(key)
}

func (ctx *Context) WriteResponse(statusCode int, resp interface{}) error {
	ctx.StopChain()

	data, err := json.Marshal(resp)
	if err != nil {
		return err // TODO: handle
	}

	if ctx.w.Header().Get("Content-Type") == "" {
		ctx.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	ctx.w.WriteHeader(statusCode)

	_, err = ctx.w.Write(data)
	if err != nil {
		return err // TODO: handle
	}

	return nil
}
