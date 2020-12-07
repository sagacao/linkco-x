package xnet

import (
	"os"

	"github.com/nothollyhigh/kiss/log"
	"github.com/nothollyhigh/kiss/net"
	"github.com/nothollyhigh/kiss/util"

	// "kisscluster/proto"
	"time"
)

// StartServer start tcp server
func StartServer(svrName string, svrAddr string, handleProto func(*net.TcpServer)) *net.TcpServer {
	tcpServer := net.NewTcpServer(svrName)

	if handleProto != nil {
		handleProto(tcpServer)
	}

	util.Go(func() {
		tcpServer.Start(svrAddr)
	})

	return tcpServer
}

// StopServer stop tcp server
func StopServer(tcpServer *net.TcpServer, svrName string) {
	tcpServer.StopWithTimeout(time.Second*5, func() {
		log.Error("%s Stop timeout", svrName)
		os.Exit(-1)
	})
}
