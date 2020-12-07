package proto

const (
	RPC_METHOD_UPDATE_SERVER_INFO = "update server info"
	RPC_METHOD_UPDATE_USER_INFO   = "OnUpdateUserInfo"

	CMD_CENTER_UPDATE_GAME_LIST_NOTIFY uint32 = 1
)

type CenterUpdateServerInfoReq struct {
	ServerInfo
}

type CenterUpdateServerInfoRsp struct {
	Code int
	Msg  string
}

type CenterUpdateUserInfoReq struct {
	SvrID  int
	UserID string
	Online int
}

type CenterUpdateUserInfoRsp struct {
	Code int
	Msg  string
}
