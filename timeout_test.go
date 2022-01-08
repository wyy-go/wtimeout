package wtimeout

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func emptySuccessResponse(c *gin.Context) {
	time.Sleep(1 * time.Second)
	c.String(http.StatusOK, "")
}

func TestTimeout(t *testing.T) {
	r := gin.New()
	r.Use(New(WithTimeout(100 * time.Microsecond)))
	r.GET("/", emptySuccessResponse)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestCallback(t *testing.T) {
	r := gin.New()
	r.Use(New(WithTimeout(100 * time.Microsecond),WithCallBack(func(request *http.Request) {
		t.Log("=========callback=========")
	})))
	r.GET("/", emptySuccessResponse)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}


func TestCustomResponse(t *testing.T) {
	r := gin.New()
	r.Use(New(WithTimeout(1 * time.Second),WithCustomMsg("custom response")))

	r.GET("/", emptySuccessResponse)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Equal(t, "custom response", w.Body.String())
}

func TestHttpCode(t *testing.T) {
	r := gin.New()
	r.Use(New(WithTimeout(100 * time.Microsecond),WithErrorHttpCode(http.StatusRequestTimeout)))
	r.GET("/", emptySuccessResponse)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestTimeout, w.Code)

}

func panicResponse(c *gin.Context) {
	panic("test")
}

func TestPanic(t *testing.T) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(New(WithTimeout(1 * time.Second)))
	r.GET("/", panicResponse)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
