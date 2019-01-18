package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"pb_test/ty"
	"pb_test/ws"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

// var addr = flag.String("addr", "192.168.1.168:8830", "http service address")

func main() {
	flag.Parse()
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt)

	var reqcount int = 100
	exitgo := make(chan struct{}, reqcount)

	var packetCount uint32 = 0

	for i := 0; i < reqcount; i++ {
		go func(i int) {
			defer func() {
				log.Printf("go[%d] exit", i)
			}()

			ws := tyws.DefaultWSClient()
			err := ws.Open(*addr, "/echo")
			if err != nil {
				return
			}

			go ws.Run()

			for {
				select {
				case <-exitgo:
					return
				default:
					atomic.AddUint32(&packetCount, 1)
					req := ty.ReqEcho{Id: packetCount, Tm: time.Now().Unix(), Data: time.Now().Format("2006-01-02 15:04:05")}
					err := ws.SendWith(0x0203, 0x0001, packetCount, &req, func(cid uint32, buf []byte) {
						var ack ty.AckEcho
						if err := proto.Unmarshal(buf, &ack); err == nil {
							log.Printf("%+v", ack)
						} else {
							log.Println(err)
						}
					})
					if err != nil {
						log.Println(err)
						return
					}
					// time.Sleep(time.Millisecond * time.Duration(200))
				}
			}
		}(i)
		time.Sleep(time.Millisecond * time.Duration(16))
	}
	<-exit
	fmt.Println("Interrupt")

	for i := 0; i < reqcount; i++ {
		exitgo <- struct{}{}
	}

	time.Sleep(time.Second * 5)
}
