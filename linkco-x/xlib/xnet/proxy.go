package xnet

import (
	"time"

	"github.com/nothollyhigh/kiss/log"
	"github.com/nothollyhigh/kiss/net"
	"github.com/nothollyhigh/kiss/timer"
	"github.com/nothollyhigh/kiss/util"
)

// StartProxyServer start rpc server
func StartProxyServer(svrName string, svrAddr string, handleProto func(*net.TcpServer)) *net.TcpServer {
	rpcServer := net.NewRpcServer(svrName)

	if handleProto != nil {
		handleProto(rpcServer)
	}

	util.Go(func() {
		rpcServer.Start(svrAddr)
	})

	return rpcServer
}

// StopProxyServer stop rpc server
func StopProxyServer(rpcServer *net.TcpServer) {
	rpcServer.StopWithTimeout(time.Second*10, nil)
}

var (
	netengine *net.TcpEngin
	rpc       *net.RpcClient
	err       error
)

func onProxyConnected(rpcb *net.RpcClient) {
	log.Info("onProxyconnect: bbbbbb addr:[%s:%d]", rpcb.Ip(), rpcb.Port())
	if rpc != nil {
		log.Info("onProxyconnect: 11111111111 addr:[%s:%d]", rpc.Ip(), rpc.Port())
	}
}

func connnect(remoteAddr string, netengine *net.TcpEngin, codec net.ICodec, onRPCConnected func(*net.RpcClient)) error {
	rpc, err = net.NewRpcClient(remoteAddr, netengine, codec, onRPCConnected)
	if err != nil {
		log.Error("connnect failed: %v", err)
		timer.AfterFunc(time.Second*5, func() {
			connnect(remoteAddr, netengine, nil, onRPCConnected)
		})
	}
	return nil
}

// StartProxy start rpc client
func StartProxy(remoteAddr string, onRPCConnected func(*net.RpcClient), handleProto func(*net.TcpEngin)) *net.RpcClient {
	netengine = net.NewTcpEngine()

	if handleProto != nil {
		handleProto(netengine)
	}

	connnect(remoteAddr, netengine, nil, onRPCConnected)

	// rpc, err = net.NewRpcClient(remoteAddr, netengine, nil, onProxyConnected)
	// if err != nil {
	// 	log.Error("NewRpcClient failed: %v", err)
	// 	panic(err)
	// }

	// log.Info("onProxyconnect: 11111111111 addr:[%s:%d]", rpc.Ip(), rpc.Port())
	return rpc
}

// StopProxy stop rpc client
func StopProxy(rpcClient *net.RpcClient) {
	util.Go(func() {
		rpcClient.Shutdown()
	})
}
