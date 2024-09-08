package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
)

type Validator[T Context] interface {
	Validate(ctx T) error
}

func POST[T Context, S Validator[T], U interface{}](pattern string, handlers ...PreparedHandlerFunc[T, S, U]) *Handler[T] {
	return &Handler[T]{
		regFn: func(r *Router, opts *HandlerOptions[T]) {
			handleWithDecode(r, "POST", opts.patternPrefix+pattern, opts, handlers...)
			handleWithDecode(r, "OPTIONS", opts.patternPrefix+pattern, opts, handlers...) // TODO: simplify
		},
	}
}

func OPTIONS[T Context](pattern string, handlers ...HandlerFunc[T]) *Handler[T] {
	return &Handler[T]{
		regFn: func(r *Router, opts *HandlerOptions[T]) {
			handle(r, "OPTIONS", opts.patternPrefix+pattern, opts, handlers...)
		},
	}
}

func GET[T Context, S Validator[T], U interface{}](pattern string, handlers ...PreparedHandlerFunc[T, S, U]) *Handler[T] {
	return &Handler[T]{
		regFn: func(r *Router, opts *HandlerOptions[T]) {
			handleWithDecode(r, "GET", opts.patternPrefix+pattern, opts, handlers...)
			handleWithDecode(r, "OPTIONS", opts.patternPrefix+pattern, opts, handlers...) // TODO: simplify
		},
	}
}

func SGET[T Context](pattern string, handlers ...HandlerFunc[T]) *Handler[T] {
	return &Handler[T]{
		regFn: func(r *Router, opts *HandlerOptions[T]) {
			handle(r, "GET", opts.patternPrefix+pattern, opts, handlers...)
		},
	}
}

func handle[T Context](r *Router, method, pattern string, opts *HandlerOptions[T], handlers ...HandlerFunc[T]) {
	patternWithMethod := method + " " + pattern

	ctxType := reflect.TypeOf([0]T{}).Elem()
	if ctxType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf(
			"\"%s\" handler ctx has not pointer type: (%s). Possible solution is to use (*%s).",
			patternWithMethod,
			ctxType.Name(),
			ctxType.Name(),
		))
	}

	r.mux.HandleFunc(patternWithMethod, func(w http.ResponseWriter, req *http.Request) {
		ctxPointer := reflect.New(ctxType.Elem())
		ctx := ctxPointer.Interface().(T)

		ctx.SetCancellableCtx(context.WithCancel(context.Background()))
		ctx.SetHTTPWriter(w)
		ctx.SetHTTPRequest(req)

		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic occurred: %v\n%s", err, debug.Stack())

				_ = ctx.WriteResponse(http.StatusInternalServerError, nil) // TODO: error json object
			}
		}()

		for i := 0; i < len(handlers); i++ {
			select {
			case <-ctx.Done():
				return
			default:
				handlers[i](ctx)
			}
		}
	})

	if method != "OPTIONS" {
		r.log.Printf("registered handler %s\n", patternWithMethod)
	}
}

func handleWithDecode[T Context, S Validator[T], U interface{}](
	r *Router,
	method, pattern string,
	opts *HandlerOptions[T],
	handlers ...PreparedHandlerFunc[T, S, U],
) {
	patternWithMethod := method + " " + pattern

	// https://arc.net/l/quote/kdxxhrfh about zero-length T array
	ctxType := reflect.TypeOf([0]T{}).Elem()
	if ctxType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf(
			"\"%s\" handler ctx has not pointer type: (%s). Possible solution is to use (*%s).",
			patternWithMethod,
			ctxType.Name(),
			ctxType.Name(),
		))
	}

	r.mux.HandleFunc(patternWithMethod, func(w http.ResponseWriter, req *http.Request) {
		ctxPointer := reflect.New(ctxType.Elem())
		ctx := ctxPointer.Interface().(T)

		ctx.SetCancellableCtx(context.WithCancel(context.Background()))
		ctx.SetHTTPWriter(w)
		ctx.SetHTTPRequest(req)

		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic occurred: %v\n%s", err, debug.Stack())

				_ = ctx.WriteResponse(http.StatusInternalServerError, nil) // TODO: error json object
			}
		}()

		for _, ph := range opts.preHandlers {
			select {
			case <-ctx.Done():
				return
			default:
				ph(ctx)
			}
		}

		decodedReq := *new(S)

		err := ctx.Decode(&decodedReq)
		if err != nil {
			if opts.errHandler != nil {
				data, statusCode := opts.errHandler(ctx, err)

				select {
				case <-ctx.Done():
					return
				default:
					err = ctx.WriteResponse(statusCode, data)
					if err != nil {
						log.Println("error writing decode error response:", err)
					}
				}
			}
			return
		}

		for _, h := range handlers {
			select {
			case <-ctx.Done():
				return
			default:
				resp, statusCode := h(ctx, &decodedReq)
				if resp != (*U)(nil) {
					err = ctx.WriteResponse(statusCode, resp)
					if err != nil {
						log.Println("error writing response:", err)
						return
					}
				}
			}
		}
	})

	if method != "OPTIONS" {
		r.log.Printf("registered handler %s\n", patternWithMethod)
	}
}
