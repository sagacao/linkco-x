package app

import (
	"context"
	"flag"
	"io"
	"linkco-x/xlib/xredis"
	"linkco-x/xutils"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/nothollyhigh/kiss/log"
	"github.com/nothollyhigh/kiss/mysql"
	"github.com/nothollyhigh/kiss/net"
)

var (
	appVersion = ""
	Config     = &AddrConfig{}
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
	confpath   = flag.String("config", "./conf/world.json", "config file path, default is conf/world.json")
	redisconf  = flag.String("redis", "./conf/redis.json", "config file path, default is conf/redis.json")
	mysqlconf  = flag.String("mysql", "./conf/mysql.json", "config file path, default is conf/mysql.json")

	logout = io.Writer(nil)
	ctx    = context.Background()

	RedisConf xredis.RedisMgrConfig
	MysqlConf mysql.MgrConfig
	TcpServer *net.TcpServer
	RpcClient *net.RpcClient
)

type AddrConfig struct {
	SvrID      int    `json:"SvrID"`
	Debug      bool   `json:"Debug"`
	LogDir     string `json:"LogDir"`
	Refresh    int    `json:"Refresh"`
	SvrAddr    string `json:"SvrAddr"`
	ProxyAddrs string `json:"ProxyAddrs"`
}

func initConfig() {
	flag.Parse()

	filename := *confpath
	err := xutils.ReadJSON(&Config, filename)
	if err != nil {
		log.Panic("initConfig Failed: %v", err)
	}

	filename = *redisconf
	RedisConf = make(xredis.RedisMgrConfig)
	err = xutils.ReadJSON(&RedisConf, filename)
	if err != nil {
		log.Panic("initConfig Failed: %v", err)
	}

	filename = *mysqlconf
	MysqlConf = make(mysql.MgrConfig)
	err = xutils.ReadJSON(&MysqlConf, filename)
	if err != nil {
		log.Panic("initConfig Failed: %v", err)
	}
}

func initLog() {
	var (
		fileWriter = &log.FileWriter{
			RootDir:     Config.LogDir + time.Now().Format("20060102/"),
			DirFormat:   "",
			FileFormat:  "20060102.log",
			MaxFileSize: 1024 * 1024 * 32,
			EnableBufio: false,
		}
	)
	if Config.Debug {
		logout = io.MultiWriter(os.Stdout, fileWriter)
	} else {
		logout = fileWriter
	}

	log.SetOutput(logout)

	configData, _ := json.MarshalIndent(Config, "", "    ")
	log.Info("config: %v\n%v", *confpath, string(configData))
}

func init() {
	initConfig()
	initLog()
}

func Run(version string) {
	appVersion = version
	log.Info("app version: '%v'", version)

}
