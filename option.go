package wtimeout

import (
	"net/http"
	"time"
)

// Option for timeout
type Option func(*Options)

type CallBackFunc func(*http.Request)

// Options struct
type Options struct {
	timeout       time.Duration
	callBack      CallBackFunc
	errorHttpCode int
	customMsg     string
}

// WithTimeout set timeout
func WithTimeout(timeout time.Duration) Option {
	return func(t *Options) {
		t.timeout = timeout
	}
}

func WithErrorHttpCode(code int) Option {
	return func(t *Options) {
		t.errorHttpCode = code
	}
}

func WithCustomMsg(s string) Option {
	return func(t *Options) {
		t.customMsg = s
	}
}

func WithCallBack(f CallBackFunc) Option {
	return func(t *Options) {
		t.callBack = f
	}
}
