package app

import (
	"linkco-x/proto"
	"time"

	"github.com/nothollyhigh/kiss/log"
)

func reportToCenter(uid string) {
	var (
		req = &proto.CenterUpdateUserInfoReq{
			SvrID:  Config.SvrID,
			UserID: uid,
			Online: 1,
		}
		rsp = &proto.CenterUpdateUserInfoRsp{}
	)

	err := RpcClient.Call(proto.RPC_METHOD_UPDATE_USER_INFO, req, rsp, time.Second*3)
	if err != nil {
		log.Error("reportToCenter updateInfo failed: %v", err)
		return
	}

	if rsp.Code == 0 {
		log.Info("reportToCenter updateInfo success")
	} else {
		log.Error("reportToCenter updateInfo failed, code: %v, msg: %v", rsp.Code, rsp.Msg)
	}
}
