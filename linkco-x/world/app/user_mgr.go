package app

import (
	"github.com/nothollyhigh/kiss/sync"

	"linkco-x/proto"
	"linkco-x/xlib/xredis"
	"linkco-x/xutils"

	"github.com/nothollyhigh/kiss/log"
	"github.com/nothollyhigh/kiss/net"
	"github.com/nothollyhigh/kiss/util"
)

var (
	userMgr = &UserMgr{
		users: map[string]*User{},
	}
)

type UserMgr struct {
	sync.RWMutex
	users map[string]*User
}

func (mgr *UserMgr) add(uid string, user *User) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.users[uid] = user
}

func (mgr *UserMgr) delete(uid string) {
	mgr.Lock()
	defer mgr.Unlock()

	delete(mgr.users, uid)
}

func (mgr *UserMgr) get(uid string) *User {
	mgr.RLock()
	defer mgr.RUnlock()

	user, ok := mgr.users[uid]
	if ok {
		return user
	}
	return nil
}

func (mgr *UserMgr) kickClient(conn *net.TcpClient, err error) {
	conn.Stop()
}

func (mgr *UserMgr) broadcast(cmd uint32, v interface{}) {
	mgr.RLock()
	defer mgr.RUnlock()

	msg := proto.NewMessage(cmd, v)

	for _, user := range mgr.users {
		user.sendMsg(msg)
	}

	log.Info("Broadcast to %d clients: %v", len(mgr.users), string(msg.Body()))
}

func (mgr *UserMgr) onDisconnect(userid string, conn *net.TcpClient) {
	conn.OnClose("disconnected", func(*net.TcpClient) {
		log.Info("user:[%s] disconnect:%v", userid, conn)
		user := mgr.get(userid)
		userMgr.delete(userid)
		user.disonnnect()
	})
}

// CMDLogin CMD--func-login
func CMDLogin(conn *net.TcpClient, msg net.IMessage) {
	var (
		err error
		req = &proto.LoginReq{}
		rsp = &proto.LoginRsp{}
	)

	if err = json.Unmarshal(msg.Body(), req); err != nil {
		rsp.Code = -1
		rsp.Msg = "invaid json"
		conn.SendMsgWithCallback(proto.NewMessage(proto.CMD_LOGIN_RSP, rsp), userMgr.kickClient)
		return
	}

	userid, err := getUserID(req.Account)
	if err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		conn.SendMsgWithCallback(proto.NewMessage(proto.CMD_LOGIN_RSP, rsp), userMgr.kickClient)
		return
	}

	user := userMgr.get(userid)
	if user == nil {
		userdata, err := xredis.Redis("user").RdbHGet(getKey(userid), userid)
		log.Info("------------------->%v %v", userdata, err)
		if err != nil {
			if err != nil {
				rsp.Code = -1
				rsp.Msg = err.Error()
				conn.SendMsgWithCallback(proto.NewMessage(proto.CMD_LOGIN_RSP, rsp), userMgr.kickClient)
				return
			}
		} else {
			if userdata == nil {
				user = NewUser(userid, req.Account, "xxx", conn)
				xredis.Redis("user").RdbHSet(getKey(userid), userid, xutils.JSONMarshalToString(user))
			} else {
				user = &User{}
				err = user.loadUser(userdata.([]byte))
				if err != nil {
					rsp.Code = -1
					rsp.Msg = err.Error()
					conn.SendMsgWithCallback(proto.NewMessage(proto.CMD_LOGIN_RSP, rsp), userMgr.kickClient)
					return
				}
			}
		}
		userMgr.add(userid, user)
		log.Info("CMDLogin connect: %v", userid)
	} else {
		userMgr.kickClient(user.MConn, nil)
		user.reconnect(conn)
		log.Info("CMDLogin reconnect: %v", userid)
	}
	userMgr.onDisconnect(userid, conn)

	util.Go(func() {
		reportToCenter(userid)
	})
	rsp.Msg = "Login Success"
	rsp.Name = userid
	conn.SendMsg(proto.NewMessage(proto.CMD_LOGIN_RSP, rsp))
}
