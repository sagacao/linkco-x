package app

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/nothollyhigh/kiss/log"
	"github.com/nothollyhigh/kiss/net"
	"github.com/nothollyhigh/kiss/util"
)

var (
	appVersion = ""
	Config     = &AddrConfig{}
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
	confpath   = flag.String("config", "./conf/center.json", "config file path, default is conf/center.json")

	logout = io.Writer(nil)

	TcpServer *net.TcpServer
	RpcClient *net.RpcClient
)

type AddrConfig struct {
	Debug      bool   `json:"Debug"`
	LogDir     string `json:"LogDir"`
	Refresh    int    `json:"Refresh"`
	SvrAddr    string `json:"SvrAddr"`
	ProxyAddrs string `json:"ProxyAddrs"`
}

func initAddrConfig() {
	flag.Parse()

	filename := *confpath
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panic("initConfig ReadFile Failed: %v", err)
	}

	data = util.TrimComment(data)
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Panic("initConfig json.Unmarshal Failed: %v", err)
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
	initAddrConfig()
	initLog()

	// svrMgr.run()
}

func Run(version string) {
	appVersion = version
	log.Info("app version: '%v'", version)

}
