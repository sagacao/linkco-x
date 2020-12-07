package app

import (
	"github.com/nothollyhigh/kiss/log"
	"github.com/nothollyhigh/kiss/sync"

	"linkco-x/proto"

	"github.com/nothollyhigh/kiss/net"
)

type userWorldInfo struct {
	SvrID  int
	UserID string
	Online int
	client *net.TcpClient
}

type userMgr struct {
	sync.RWMutex
	users map[string]*userWorldInfo
}

func OnUpdateUserInfo(ctx *net.RpcContext) {
	var (
		err error
		// code int
		req = &proto.CenterUpdateUserInfoReq{}
		rsp = &proto.CenterUpdateUserInfoRsp{}
	)

	if err = ctx.Bind(req); err != nil {
		return
	}
	log.Info("OnUpdateUserInfo: %v", req.UserID)

	rsp.Code = 0
	rsp.Msg = "success"
	ctx.Write(rsp)
	// ctx.Client()

	// svr := &ServerInfo{req.ServerInfo, ctx.Client()}
	// code, err = svrMgr.Add(svr)
	// if err != nil {
	// 	rsp.Code = code
	// 	rsp.Msg = err.Error()
	// } else {
	// 	ctx.Client().OnClose("DeleServer", func(c *net.TcpClient) {
	// 		svrMgr.Delete(svr)
	// 	})
	// }

	// ctx.Write(rsp)

	// log.Info("onUpdateServerInfo: %v", string(ctx.Body()))
}
