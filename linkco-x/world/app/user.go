package app

import (
	"linkco-x/xlib/xredis"
	"strings"

	"github.com/nothollyhigh/kiss/log"
	"github.com/nothollyhigh/kiss/net"
	uuid "github.com/satori/go.uuid"
)

var ()

type User struct {
	MConn    *net.TcpClient
	MAccount string `json:"account"`
	MUID     string `json:"id"`
	MName    string `json:"name"`
}

func NewUser(uid, account, name string, conn *net.TcpClient) *User {
	return &User{
		MConn:    conn,
		MAccount: account,
		MUID:     uid,
		MName:    name,
	}
}

func getUserID(account string) (string, error) {
	userid, err := xredis.Redis("account").RdbGet(account)
	if err == nil {
		return userid, nil
	}

	// log.Error("getUserID account:[%s] err:%v %v", account, err.Error(), strings.Contains(err.Error(), "nil returned"))
	if strings.Contains(err.Error(), "nil returned") {
		return uuid.NewV4().String(), nil
	}
	log.Error("getUserID account:[%s] err:%v", account, err)
	return "", err
}

func getKey(userid string) string {
	return userid
}

func (u *User) reconnect(conn *net.TcpClient) {
	u.MConn = conn
}

func (u *User) sendMsg(msg *net.Message) {
	u.MConn.SendMsg(msg)
}

func (u *User) loadUser(userdata []byte) error {
	err := json.Unmarshal(userdata, &u)
	if err != nil {
		log.Error("json.Unmarshal Failed: %v", err)
		return err
	}
	return nil
}

func (u *User) disonnnect() {

}
