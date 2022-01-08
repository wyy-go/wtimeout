package wtimeout

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// TimeoutWriter is a writer with memory buffer
type TimeoutWriter struct {
	gin.ResponseWriter
	body         *bytes.Buffer
	headers      http.Header
	mu           sync.Mutex
	timeout      bool
	wroteHeaders bool
	code         int
}

// newTimeoutWriter will return a timeout.Writer pointer
func newTimeoutWriter(w gin.ResponseWriter, buf *bytes.Buffer) *TimeoutWriter {
	return &TimeoutWriter{ResponseWriter: w, body: buf, headers: make(http.Header)}
}

func (tw *TimeoutWriter) Timeout() {
	tw.timeout = true
}

// Write will write data to response body
func (tw *TimeoutWriter) Write(data []byte) (int, error) {
	if tw.timeout {
		return 0, nil
	}
	tw.mu.Lock()
	defer tw.mu.Unlock()

	return tw.body.Write(data)
}

// WriteString will write string to response body
func (tw *TimeoutWriter) WriteString(s string) (int, error) {
	return tw.Write([]byte(s))
}

// WriteHeader will write http status code
func (tw *TimeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)
	if tw.timeout {
		return
	}
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.writeHeader(code)
}

func (tw *TimeoutWriter) writeHeader(code int) {
	tw.wroteHeaders = true
	tw.code = code
}

// Header will get response headers
func (tw *TimeoutWriter) Header() http.Header {
	return tw.headers
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid http status code: %d", code))
	}
}

func (tw *TimeoutWriter) WriteHeaderNow() {}
