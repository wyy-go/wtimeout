# Timeout

![GitHub Repo stars](https://img.shields.io/github/stars/wyy-go/wtimeout?style=social)
![GitHub](https://img.shields.io/github/license/wyy-go/wtimeout)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/wyy-go/wtimeout)
![GitHub CI Status](https://img.shields.io/github/workflow/status/wyy-go/wtimeout/ci?label=CI)
[![Go Report Card](https://goreportcard.com/badge/github.com/wyy-go/wtimeout)](https://goreportcard.com/report/github.com/wyy-go/wtimeout)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/wyy-go/wtimeout?tab=doc)
[![codecov](https://codecov.io/gh/wyy-go/wtimeout/branch/main/graph/badge.svg)](https://codecov.io/gh/wyy-go/wtimeout)



Timeout wraps a handler and aborts the process of the handler if the timeout is reached.

## Example

```go
package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wyy-go/wtimeout"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func AccessLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("[start]AccessLog")
		ctx.Next()
		log.Println("[end]AccessLog")
	}
}


func main() {

	// create new gin without any middleware
	engine := gin.Default()

	customMsg := `{"code": -1, "msg":"http: Handler timeout"}`
	// add timeout middleware with 2 second duration
	engine.Use(wtimeout.New(
		wtimeout.WithTimeout(2*time.Second),
		wtimeout.WithErrorHttpCode(http.StatusRequestTimeout), // optional
		wtimeout.WithCustomMsg(customMsg),                   // optional
		wtimeout.WithCallBack(func(r *http.Request) {
			fmt.Println("timeout happen, url:", r.URL.String())
		}), // optional
	))
	// create a handler that will last 1 seconds
	engine.GET("/short", short)

	// create a handler that will last 5 seconds
	engine.GET("/long", AccessLog(), long)

	// create a handler that will last 5 seconds but can be canceled.
	engine.GET("/long2", long2)

	// create a handler that will last 20 seconds but can be canceled.
	engine.GET("/long3", long3)

	engine.GET("/boundary", boundary)

	// run the server
	log.Fatal(engine.Run(":8080"))
}

func short(c *gin.Context) {
	time.Sleep(1 * time.Second)
	c.JSON(http.StatusOK, gin.H{"hello": "short"})
}

func long(c *gin.Context) {
	fmt.Println("handler-long1, do something...")
	time.Sleep(3 * time.Second)
	fmt.Println("handler-long2, do something...")
	time.Sleep(3 * time.Second)
	fmt.Println("handler-long3, do something...")
	c.JSON(http.StatusOK, gin.H{"hello": "long"})
}

func boundary(c *gin.Context) {
	time.Sleep(2 * time.Second)
	c.JSON(http.StatusOK, gin.H{"hello": "boundary"})
}

func long2(c *gin.Context) {
	if doSomething(c.Request.Context()) {
		c.JSON(http.StatusOK, gin.H{"hello": "long2"})
	}
}

func long3(c *gin.Context) {
	// request a slow service
	// see  https://github.com/vearne/gin-timeout/blob/master/example/slow_service.go
	url := "http://localhost:8882/hello"
	// Notice:
	// Please use c.Request.Context(), the handler will be canceled where timeout event happen.
	req, _ := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, url, nil)
	client := http.Client{Timeout: 100 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// Where timeout event happen, a error will be received.
		fmt.Println("error1:", err)
		return
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error2:", err)
		return
	}
	fmt.Println(s)
}

// A cancelCtx can be canceled.
// When canceled, it also cancels any children that implement canceler.
func doSomething(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		fmt.Println("doSomething is canceled.")
		return false
	case <-time.After(5 * time.Second):
		fmt.Println("doSomething is done.")
		return true
	}
}
```
