package xnet

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nothollyhigh/kiss/log"
	"github.com/nothollyhigh/kiss/net"
	"github.com/nothollyhigh/kiss/util"
)

// example:
// func registerRouter(router *gin.Engine) {
// 	router.GET("/hello", func(c *gin.Context) {
// 		log.Info("onHello")
// 		c.String(http.StatusOK, "hello")
// 	})
// }

// NewWebServer("Web", ":8080", registerRouter, true, time.Second*5, nil, "", "")

// NewWebServer new web server
func NewWebServer(tag, addr string, routerHandle func(*gin.Engine), isDebug bool, timeout time.Duration, opt *net.SocketOpt, certFile string, keyFile string) {
	if isDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	if routerHandle != nil {
		routerHandle(router)
	}

	// net.ServeHttp(tag, addr, router, time.Second*5, nil)

	svr, err := net.NewHttpServer(tag, addr, router, timeout, opt, func() {
		os.Exit(0)
	})
	if err != nil {
		log.Fatal("[HttpServer %v]: ServeTLS failed: %v", tag, err)
	} else {
		util.Go(func() {
			if certFile != "" && keyFile != "" {
				svr.ServeTLS(certFile, keyFile)
			} else {
				svr.Serve()
			}
		})

		if isDebug {
			svr.EnablePProf(tag)
		}
	}

	util.HandleSignal(func(sig os.Signal) {
		if sig == syscall.SIGTERM || sig == syscall.SIGINT {
			svr.Shutdown()
			os.Exit(0)
		}
	})
}

// Get 发送GET请求
// url：         请求地址
// response：    请求返回的内容
func Get(url string) string {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return result.String()
}

// Post 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如：application/json
// content：     请求放回的内容
func Post(url string, data interface{}, contentType string) string {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}
