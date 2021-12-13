package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	log "github.com/store_server/logger"
	"github.com/store_server/store_server_rpc/g"
	rpc "github.com/store_server/store_server_rpc/rpc/rpcServer"
	"github.com/store_server/utils"
	"github.com/store_server/utils/common"
)

var (
	GitTag   = "tag"
	Version  = "dev"
	Build    = "2020-09-09"
	serverIP = "127.0.0.1"
)

func pv(code int) {
	fmt.Fprintf(os.Stdout, "GitTag: %s\n", GitTag)
	fmt.Fprintf(os.Stdout, "Version: %s\n", Version)
	fmt.Fprintf(os.Stdout, "Build: %s\n", Build)
	os.Exit(code)
}

func flushSentry() {
	tk := time.NewTicker(3 * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-tk.C:
			sentry.Flush(10 * time.Second)
		}
	}
}

func main() {
	version := flag.Bool("v", false, "version")
	path := flag.String("c", "", "config path")
	debug := flag.Bool("d", false, "debug")
	logPath := flag.String("l", "/data/apps/store_server_rpc/logs/store_server_rpc.log", "log path")
	level := flag.String("log_level", "debug", "log level")
	flag.Parse()
	if *version {
		pv(0)
	}
	if *path == "" {
		pv(1)
	}

	g.InitConf(*path, *debug)
	log.Init("store_server_rpc", *logPath, g.Config().Rpc.Debug) //初始化logger
	log.InitStructLog(*level, *logPath, "store_server")

	serverIP, _ = common.GetIPByMultiAddr([]string{"eth1", "eth0", "en0"})
	go flushSentry()
	instance, err := rpc.NewDefaultRpcServer()
	if err != nil {
		log.Fatal(err)
	}
	s := utils.NewService(instance)
	s.Forever()
	sentry.CaptureException(fmt.Errorf("%v server shutdown", os.Getpid()))
	sentry.Flush(time.Second * 1)
}
