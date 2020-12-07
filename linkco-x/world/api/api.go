package api

import (
	"github.com/nothollyhigh/kiss/net"

	"linkco-x/proto"
	"linkco-x/world/app"
)

// ServerCmdRegister register server cmd
func ServerCmdRegister(svr *net.TcpServer) {
	// svr.Handle(proto.CMD_CENTER_UPDATE_GAME_LIST_NOTIFY, onUpdateServerInfo)
	svr.Handle(proto.CMD_LOGIN_REQ, app.CMDLogin)
}

// ClientCmdRegister register client cmd
func ClientCmdRegister(svr *net.TcpEngin) {
	// svr.Handle(proto.CMD_CENTER_UPDATE_GAME_LIST_NOTIFY, onUpdateServerInfo)
}
