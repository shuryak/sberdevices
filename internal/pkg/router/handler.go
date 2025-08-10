package router

type Handler[T Context] struct {
	opts  HandlerOptions[T]
	regFn regFn[T]
}

type regFn[T Context] func(r *Router, opts *HandlerOptions[T])

type HandlerOptions[T Context] struct {
	errHandler    ErrHandlerFunc[T]
	preHandlers   []HandlerFunc[T]
	patternPrefix string
}

type HandlerFunc[T Context] func(ctx T)

type PreparedHandlerFunc[T Context, S Validator[T], U interface{}] func(ctx T, req *S) (*U, int)

type ErrHandlerFunc[T Context] func(ctx T, err error) (interface{}, int)

func (h *Handler[T]) SetErrHandler(errHandler ErrHandlerFunc[T]) *Handler[T] {
	h.opts.errHandler = errHandler
	return h
}

func (h *Handler[T]) SetPreHandler(preHandlers ...HandlerFunc[T]) *Handler[T] {
	h.opts.preHandlers = preHandlers
	return h
}

func (h *Handler[T]) register(r *Router) {
	h.regFn(r, &h.opts)
}
