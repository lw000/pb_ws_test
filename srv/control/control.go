package control

import (
	"demo/pb_test/ty"
	"sync"
	"time"
	"tuyue/tuyue_common/network/ws/packet"
	"tuyue/tuyue_common/network/ws/srv/hub"
	"tuyue/tuyue_common/schedule"
	"tuyue/tuyue_common/utilty"

	log "github.com/alecthomas/log4go"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var (
	hub  *tyhub.Hub
	sche *tyschedule.Schedule
	lock sync.Mutex
)

const (
	reply_millisecond = 200
)

func init() {
	registerHub()
}

func AckMessage(c *websocket.Conn, data []byte) {
	//rand.Intn(reply_millisecond)
	if err := c.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Error(err)
	}
	//cdata := make(map[interface{}]interface{})
	//cdata["c"] = c
	//cdata["mt"] = mt
	//cdata["data"] = data
	//sche.AddTask(1, cdata)
}

func AddData() {
	for {
		cdata := make(map[interface{}]interface{})
		cdata["c"] = tyutilty.RandomIntger(8)
		cdata["mt"] = 1
		cdata["data"] = tyutilty.UUID()
		sche.AddTask(1, cdata)

		time.Sleep(time.Second * time.Duration(1))
	}
}

func registerHub() {
	sche = tyschedule.NewSchedule()
	err := sche.Start(func(taskId, data interface{}) {
		m := data.(map[interface{}]interface{})
		c := m["c"].(*websocket.Conn)
		mt := m["mt"].(int)
		cdata := m["data"].([]byte)
		lock.Lock()
		if err := c.WriteMessage(mt, cdata); err != nil {
			log.Error(err)
		}
		lock.Unlock()
	})

	if err != nil {

	}

	//go AddData()

	hub = tyhub.NewWSHub()
	hub.Handle(0x0203, 0x0001, func(c *websocket.Conn, pk *typacket.Packet) {
		var req ty.ReqEcho
		if err := proto.Unmarshal(pk.Data(), &req); err != nil {
			log.Error(err)
			return
		}
		log.Info("[%s] [%d] %+v", c.RemoteAddr().String(), pk.RequestId(), req)

		ack := typacket.NewPacket(pk.Mid(), pk.Sid(), pk.RequestId())
		_ = ack.EncodeProto(&ty.AckEcho{Code: req.Id, Data: req.Data})
		AckMessage(c, ack.Data())
	})

	hub.Handle(0x0203, 0x0002, func(c *websocket.Conn, pk *typacket.Packet) {
		var req ty.ReqLogin
		if err := proto.Unmarshal(pk.Data(), &req); err != nil {
			log.Error(err)
			return
		}
		log.Info("[%s] [%d] %+v", c.RemoteAddr().String(), pk.RequestId(), req)
		ack := typacket.NewPacket(pk.Mid(), pk.Sid(), pk.RequestId())
		_ = ack.EncodeProto(&ty.AckLogin{Code: req.Id, Data: req.Data})
		AckMessage(c, ack.Data())
	})

	hub.Handle(0x0203, 0x0003, func(c *websocket.Conn, pk *typacket.Packet) {
		var req ty.ReqLogout
		if err := proto.Unmarshal(pk.Data(), &req); err != nil {
			log.Error(err)
			return
		}

		log.Info("[%s] [%d] %+v", c.RemoteAddr().String(), pk.RequestId(), req)

		ack := typacket.NewPacket(pk.Mid(), pk.Sid(), pk.RequestId())
		_ = ack.EncodeProto(&ty.AckLogout{Code: req.Id, Data: req.Data})
		AckMessage(c, ack.Data())
	})
}

func GetHub() *tyhub.Hub {
	return hub
}
