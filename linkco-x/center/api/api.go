package api

import (
	"github.com/nothollyhigh/kiss/net"

	"linkco-x/center/app"
	"linkco-x/proto"
)

// ServerCmdRegister register server cmd
func ServerCmdRegister(svr *net.TcpServer) {
	// svr.Handle(proto.CMD_CENTER_UPDATE_GAME_LIST_NOTIFY, onUpdateServerInfo)
	svr.HandleRpcMethod(proto.RPC_METHOD_UPDATE_SERVER_INFO, app.OnUpdateServerInfo)
	svr.HandleRpcMethod(proto.RPC_METHOD_UPDATE_USER_INFO, app.OnUpdateUserInfo)
}

// ClientCmdRegister register client cmd
func ClientCmdRegister(svr *net.TcpEngin) {
	// svr.Handle(proto.CMD_CENTER_UPDATE_GAME_LIST_NOTIFY, onUpdateServerInfo)
}
