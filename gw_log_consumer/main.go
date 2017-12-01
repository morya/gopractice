package main

import (
	"os"
	"os/signal"
	"syscall"

	"flag"

	"github.com/morya/utils/log"
)

var (
	flagRedisAddr    = flag.String("redisaddr", "127.0.0.1:6379", "redis address")
	flagRedisKey     = flag.String("rediskey", "filebeat", "redis keyname for pull data")
	flagUnixSockAddr = flag.String("unixsocket", "", "unix socket address")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var app = NewApp(*flagRedisAddr, *flagRedisKey)
	if !app.init(*flagUnixSockAddr, `ELK_DATA (.*) ELK_END`) {
		return
	}

	var onExit = make(chan os.Signal)
	signal.Notify(onExit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		var sig = <-onExit
		log.Infof("on signal [%v], prepare to exit...", sig.String())
		app.cancel()
	}()

	app.Run()

	log.Info("Bye")
}
