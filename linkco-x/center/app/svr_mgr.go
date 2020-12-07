package app

import (
	"github.com/nothollyhigh/kiss/log"
	"github.com/nothollyhigh/kiss/net"

	"linkco-x/proto"
)

func OnUpdateServerInfo(ctx *net.RpcContext) {
	var (
		err error
		// code int
		req = &proto.CenterUpdateServerInfoReq{}
		rsp = &proto.CenterUpdateServerInfoRsp{}
	)

	if err = ctx.Bind(req); err != nil {
		rsp.Code = -1
		rsp.Msg = "invalid body"
		ctx.Write(rsp)
		return
	}

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

	ctx.Write(rsp)

	log.Info("onUpdateServerInfo: %v", string(ctx.Body()))
}
