package main

import (
	"linkco-x/world/api"
	"linkco-x/world/app"
	"linkco-x/xlib/xnet"
	"linkco-x/xlib/xredis"
	"os"
	"syscall"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/nothollyhigh/kiss/log"
	"github.com/nothollyhigh/kiss/net"
	"github.com/nothollyhigh/kiss/util"
)

var version = "0.0.1"
var svrName = "World"

func stop() {
	ch := make(chan int, 1)

	util.Go(func() {
		xredis.RedisFinal()
		if app.RpcClient != nil {
			app.RpcClient.Shutdown()
		}

		if app.TcpServer != nil {
			xnet.StopServer(app.TcpServer, svrName)
		}
		ch <- 1
	})

	select {
	case <-ch:
	case <-time.After(time.Second * 5):
		log.Error("%s Stop timeout", svrName)
	}

	log.Info("%s stop # %v #", svrName, time.Now().Format("2006-01-02 15:04"))
}

func onProxyconnect(rpc *net.RpcClient) {
	log.Info("onProxyconnect: addr:[%s:%d]", rpc.Ip(), rpc.Port())
	app.RpcClient = rpc
}

func start() {
	app.TcpServer = xnet.StartServer(svrName, app.Config.SvrAddr, api.ServerCmdRegister)

	if app.RedisConf != nil {
		xredis.RedisInit(app.RedisConf)
		xredis.ForEach(func(tag string, rds redis.Conn) {
			if rds != nil {
				log.Info("redis:[%s] success to conn", tag)
			} else {
				log.Info("redis:[%s] fail to conn", tag)
			}
		})
	}

	// if app.MysqlConf != nil {
	// 	xdb.DBConnInit(app.MysqlConf)
	// }

	if app.Config.ProxyAddrs != "" {
		app.RpcClient = xnet.StartProxy(app.Config.ProxyAddrs, onProxyconnect, api.ClientCmdRegister)
	}

	log.Info("%s start # %v #", svrName, time.Now().Format("2006-01-02 15:04"))
}

func main() {
	start()
	app.Run(version)

	util.HandleSignal(func(sig os.Signal) {
		if sig == syscall.SIGTERM || sig == syscall.SIGINT {
			stop()
			os.Exit(0)
		}
	})
}
