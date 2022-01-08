package wtimeout

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants/v2"
)

const (
	defaultTimeout = 5 * time.Second
)

// New wraps a handler and aborts the process of the handler if the timeout is reached
func New(opts ...Option) gin.HandlerFunc {
	options := Options{
		timeout:       defaultTimeout,
		callBack:      nil,
		customMsg:     `{"code": -1, "msg":"http: Handler timeout"}`,
		errorHttpCode: http.StatusServiceUnavailable,
	}

	// Loop through each option
	for _, opt := range opts {
		opt(&options)
	}

	return func(c *gin.Context) {

		cp := *c //nolint: govet
		c.Abort()
		c.Keys = nil

		// sync.Pool
		buffer := defaultBufPool.GetBuff()
		tw := newTimeoutWriter(cp.Writer, buffer)

		cp.Writer = tw

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(cp.Request.Context(), options.timeout)
		defer cancel()

		cp.Request = cp.Request.WithContext(ctx)

		// Channel capacity must be greater than 0.
		// Otherwise, if the parent coroutine quit due to timeout,
		// the child coroutine may never be able to quit.
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		_ = ants.Submit(func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			cp.Next()
			finish <- struct{}{}
		})

		select {
		case p := <-panicChan:
			panic(p)

		case <-finish:
			tw.mu.Lock()
			defer tw.mu.Unlock()
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}

			if !tw.wroteHeaders {
				tw.code = http.StatusOK
			}

			tw.ResponseWriter.WriteHeader(tw.code)
			if _, err := tw.ResponseWriter.Write(buffer.Bytes()); err != nil {
				panic(err)
			}

			defaultBufPool.PutBuff(buffer)

		case <-ctx.Done():
			tw.mu.Lock()
			defer tw.mu.Unlock()

			tw.Timeout()
			tw.ResponseWriter.WriteHeader(options.errorHttpCode)
			if _, err := tw.ResponseWriter.WriteString(options.customMsg); err != nil {
				panic(err)
			}

			cp.Abort()

			// execute callback func
			if options.callBack != nil {
				options.callBack(cp.Request.Clone(context.Background()))
			}

			defaultBufPool.PutBuff(buffer)
		}
	}
}
