package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wyy-go/wtimeout"
	"log"
	"net/http"
	"time"
)

type errResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func MyRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("unknow error:%v\n", p)
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					errResponse{Code: -1, Msg: fmt.Sprintf("unknow error:%v", p)})
				return
			}
		}()
		c.Next()
	}
}

func main() {

	router := gin.Default()
	// In order to maintain flexibility,
	// you should define your own recovery middleware
	router.Use(MyRecovery())
	customMsg := `{"code": -1, "msg":"http: Handler timeout"}`
	router.Use(wtimeout.New(wtimeout.WithTimeout(3*time.Second),
		wtimeout.WithCustomMsg(customMsg)))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, time.Now().String())
	})
	router.GET("/panic", func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		x := 0
		fmt.Println(100 / x)
	})
	log.Fatal(router.Run(":8080"))
}