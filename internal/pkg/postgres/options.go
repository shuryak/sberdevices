package postgres

import "time"

type Options struct {
	maxPoolSize  int32
	connAttempts int32
	connTimeout  time.Duration
}

func NewOptions() *Options {
	return &Options{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}
}

func (o *Options) SetMaxPoolSize(maxPoolSize int32) *Options {
	o.maxPoolSize = maxPoolSize
	return o
}

func (o *Options) SetConnAttempts(connAttempts int32) *Options {
	o.connAttempts = connAttempts
	return o
}

func (o *Options) SetConnTimeout(connTimeout time.Duration) *Options {
	o.connTimeout = connTimeout
	return o
}

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 5
	_defaultConnTimeout  = time.Second
)
