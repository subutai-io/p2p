package main

import (
	"testing"

	"github.com/golang/protobuf/proto"
	ptp "github.com/subutai-io/p2p"
	"github.com/subutai-io/p2p/protocol"
)

func TestRouteData(t *testing.T) {

	dataStr := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	dataArr := []string{}
	i := 0
	for i < 99 {
		dataArr = append(dataArr, dataStr)
		i++
	}

	data := []protocol.DHTPacket{
		protocol.DHTPacket{},
		protocol.DHTPacket{
			Type: protocol.DHTPacketType_Ping,
		},
		protocol.DHTPacket{
			Type:      protocol.DHTPacketType_Ping,
			Arguments: dataArr,
		},
	}

	router := new(DHTRouter)
	for _, d := range data {
		b, _ := proto.Marshal(&d)
		if len(b) > ptp.DHTBufferSize {
			b = b[:ptp.DHTBufferSize]
		}
		router.routeData(b)
	}
}
