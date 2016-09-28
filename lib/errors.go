package ptp

import (
	"errors"
)

// ErrorType is a type of an error
type ErrorType string

var (
	// ErrorList stores known errors
	ErrorList map[ErrorType]error
)

// Types of errors
const (
	ErrUnknownError       ErrorType = "unknownerror"
	ErrIncopatibleVersion ErrorType = "unsupported"
	ErrMalformedHandshake ErrorType = "badhandshake"
	ErrPortParseFailed    ErrorType = "badport"
	ErrBadUDPAddr         ErrorType = "badudpaddr"
	ErrBadIDReceived      ErrorType = "badid"
	ErrBadDHCPData        ErrorType = "baddhcp"
)

// TError -Struct for errors
type TError struct {
	Type ErrorType
}

// InitErrors populates ErrorList with error types
func InitErrors() {
	ErrorList = make(map[ErrorType]error)
	ErrorList[ErrIncopatibleVersion] = errors.New("DHT received incompatible packet")
	ErrorList[ErrMalformedHandshake] = errors.New("DHT received malformed handshake")
	ErrorList[ErrPortParseFailed] = errors.New("DHT failed to extract port from handshake")
	ErrorList[ErrBadUDPAddr] = errors.New("DHT failed to extract UDP address from handshake")
	ErrorList[ErrBadIDReceived] = errors.New("DHT received invalid ID from client")
	ErrorList[ErrBadDHCPData] = errors.New("DHT failed to parse provided DHCP packet")
}
