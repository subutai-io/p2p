package ptp

import (
	"errors"
)

type ErrorType string

var (
	ErrorList map[ErrorType]error
)

const (
	ERR_UNKNOWN_ERROR       ErrorType = "unknownerror"
	ERR_INCOPATIBLE_VERSION ErrorType = "unsupported"
	ERR_MALFORMED_HANDSHAKE ErrorType = "badhandshake"
	ERR_PORT_PARSE_FAILED   ErrorType = "badport"
	ERR_BAD_UDP_ADDR        ErrorType = "badudpaddr"
	ERR_BAD_ID_RECEIVED     ErrorType = "badid"
	ERR_BAD_DHCP_DATA       ErrorType = "baddhcp"
)

type Error struct {
	Type ErrorType
}

func InitErrors() {
	ErrorList = make(map[ErrorType]error)
	ErrorList[ERR_INCOPATIBLE_VERSION] = errors.New("DHT received incompatible packet")
	ErrorList[ERR_MALFORMED_HANDSHAKE] = errors.New("DHT received malformed handshake")
	ErrorList[ERR_PORT_PARSE_FAILED] = errors.New("DHT failed to extract port from handshake")
	ErrorList[ERR_BAD_UDP_ADDR] = errors.New("DHT failed to extract UDP address from handshake")
	ErrorList[ERR_BAD_ID_RECEIVED] = errors.New("DHT received invalid ID from client")
	ErrorList[ERR_BAD_DHCP_DATA] = errors.New("DHT failed to parse provided DHCP packet")
}
