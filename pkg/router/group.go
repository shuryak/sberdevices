package router

func NewGroup[T Context](patternPrefix string, handlers ...*Handler[T]) *Handler[T] {
	return &Handler[T]{
		regFn: func(r *Router, opts *HandlerOptions[T]) {
			for i := 0; i < len(handlers); i++ {
				if handlers[i].opts.errHandler == nil {
					handlers[i].opts.errHandler = opts.errHandler
				}
				if handlers[i].opts.preHandlers == nil {
					handlers[i].opts.preHandlers = opts.preHandlers
				}

				handlers[i].opts.patternPrefix = opts.patternPrefix + patternPrefix
				handlers[i].register(r)
			}
		},
	}
}
