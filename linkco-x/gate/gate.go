package main

import (
	"linkco-x/gate/app"
	"os"
	"syscall"

	"github.com/nothollyhigh/kiss/util"
)

var Version = ""

func main() {
	app.Run(Version)

	util.HandleSignal(func(sig os.Signal) {
		if sig == syscall.SIGTERM || sig == syscall.SIGINT {
			app.Stop()
			os.Exit(0)
		}
	})
}
