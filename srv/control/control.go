package control

import (
	"demo/pb_test/ty"
	"tuyue/tuyue_common/network/ws/packet"
	"tuyue/tuyue_common/network/ws/srv/hub"

	log "github.com/alecthomas/log4go"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var (
	hub *tyhub.Hub
)

func init() {
	registerHub()
}

func AckMessage(c *websocket.Conn, data []byte) {
	if err := c.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Error(err)
	}
}

func registerHub() {

	hub = tyhub.NewHub()
	hub.AddHandle(0x0203, 0x0001, func(c *websocket.Conn, pk *typacket.Packet) {
		var req ty.ReqEcho
		if err := proto.Unmarshal(pk.Data(), &req); err != nil {
			log.Error(err)
			return
		}
		log.Info("[%s] [%d] %+v", c.RemoteAddr().String(), pk.ClientId(), req)

		ack := typacket.NewPacket(pk.Mid(), pk.Sid(), pk.ClientId())
		data, er := proto.Marshal(&ty.AckEcho{Code: req.GetId(), Data: req.Data})
		if er != nil {
			return
		}
		_ = ack.Encode(data)
		AckMessage(c, ack.Data())
	})

	hub.AddHandle(0x0203, 0x0002, func(c *websocket.Conn, pk *typacket.Packet) {
		var req ty.ReqLogin
		if err := proto.Unmarshal(pk.Data(), &req); err != nil {
			log.Error(err)
			return
		}
		log.Info("[%s] [%d] %+v", c.RemoteAddr().String(), pk.ClientId(), req)
		ack := typacket.NewPacket(pk.Mid(), pk.Sid(), pk.ClientId())
		data, er := proto.Marshal(&ty.AckLogin{Code: req.Id, Data: req.Data})
		if er != nil {
			return
		}
		_ = ack.Encode(data)
		AckMessage(c, ack.Data())
	})

	hub.AddHandle(0x0203, 0x0003, func(c *websocket.Conn, pk *typacket.Packet) {
		var req ty.ReqLogout
		if err := proto.Unmarshal(pk.Data(), &req); err != nil {
			log.Error(err)
			return
		}

		log.Info("[%s] [%d] %+v", c.RemoteAddr().String(), pk.ClientId(), req)

		ack := typacket.NewPacket(pk.Mid(), pk.Sid(), pk.ClientId())
		data, er := proto.Marshal(&ty.AckLogout{Code: req.GetId(), Data: req.Data})
		if er != nil {
			return
		}
		_ = ack.Encode(data)
		AckMessage(c, ack.Data())
	})
}

func GetHub() *tyhub.Hub {
	return hub
}
