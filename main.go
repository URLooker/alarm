package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/urlooker/alarm/backend"
	"github.com/urlooker/alarm/cron"
	"github.com/urlooker/alarm/g"
	"github.com/urlooker/alarm/judge"
	"github.com/urlooker/alarm/receiver"
	"github.com/urlooker/alarm/sender"
)

func prepare() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func init() {
	prepare()

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	handleVersion(*version)
	handleHelp(*help)
	handleConfig(*cfg)

	g.InitRedisConnPool()
	judge.InitHistoryBigMap()
	sender.Init()
	backend.InitClients(g.Config.Web.Addrs)
}

func main() {
	go cron.ReadEvent()
	go cron.SyncStrategies()
	go sender.PopAllMail(g.Config.Queue.Mail)
	receiver.Start()
	log.Println("ok")
}

func handleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
}

func handleHelp(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func handleConfig(configFile string) {
	err := g.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
