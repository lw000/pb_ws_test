package main

import (
	"demo/pb_ws_test/cli/config"
	"demo/pb_ws_test/ty"
	"flag"
	"github.com/golang/protobuf/proto"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
	"tuyue/tuyue_common/ws/cli"

	log "github.com/alecthomas/log4go"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

// var addr = flag.String("addr", "192.168.1.168:8830", "http service address")

var ws *tyws.WSClient
var cfg *config.Config

var reqId int64

func SendTest() {
	var i int64
	for i = 0; i < cfg.Cfg.Count; i++ {
		go func() {
			{
				req := ty.ReqEcho{Id: atomic.AddInt64(&reqId, 1), Tm: time.Now().Unix(), Data: time.Now().Format("2006-01-02 15:04:05")}
				_ = ws.AsynSendMessage(0x0203, 0x0001, &req, func(buf []byte) {
					var ack ty.AckEcho
					if err := proto.Unmarshal(buf, &ack); err != nil {
						log.Error(err)
						return
					}
					log.Info("%+v", ack)
				})
			}

			{
				req := ty.ReqLogin{Id: atomic.AddInt64(&reqId, 1), Tm: time.Now().Unix(), Data: "login"}
				_ = ws.AsynSendMessage(0x0203, 0x0002, &req, func(buf []byte) {
					var ack ty.AckLogin
					if err := proto.Unmarshal(buf, &ack); err != nil {
						log.Error(err)
						return
					}
					log.Info("%+v", ack)
				})
			}

			{
				req := ty.ReqLogout{Id: atomic.AddInt64(&reqId, 1), Tm: time.Now().Unix(), Data: "logout"}
				_ = ws.AsynSendMessage(0x0203, 0x0003, &req, func(buf []byte) {
					var ack ty.AckLogout
					if err := proto.Unmarshal(buf, &ack); err != nil {
						log.Error(err)
						return
					}
					log.Info("%+v", ack)
				})
			}

			time.Sleep(time.Millisecond * time.Duration(16))
		}()
	}
}

func main() {
	flag.Parse()

	log.LoadConfiguration("../configs/log4go.xml")

	cfg = config.NewConfig()
	if err := cfg.Load("./conf/conf.json"); err != nil {
		return
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var err error
	ws = tyws.DefaultAyncClient(30, 1024)
	err = ws.Open("ws", *addr, "/ws")
	if err != nil {
		return
	}

	err = ws.Run()
	if err != nil {
		return
	}

	SendTest()

	for {
		select {
		case <-interrupt:
			ws.Stop()

			select {
			case <-time.After(time.Second):
			}
			return
		}
	}
}
