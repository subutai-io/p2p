package main

import (
	"testing"

	"github.com/golang/protobuf/proto"
	ptp "github.com/subutai-io/p2p/lib"
)

func TestRouteData(t *testing.T) {

	dataStr := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	dataArr := []string{}
	i := 0
	for i < 99 {
		dataArr = append(dataArr, dataStr)
		i++
	}

	data := []ptp.DHTPacket{
		ptp.DHTPacket{},
		ptp.DHTPacket{
			Type: ptp.DHTPacketType_Ping,
		},
		ptp.DHTPacket{
			Type:      ptp.DHTPacketType_Ping,
			Arguments: dataArr,
		},
	}

	router := new(DHTRouter)
	for _, d := range data {
		b, _ := proto.Marshal(&d)
		if len(b) > 1024 {
			b = b[:1024]
		}
		router.routeData(b)
	}
}
