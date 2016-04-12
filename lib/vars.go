package ptp

import (
	"time"
)

const PACKET_VERSION string = "4"

var SUPPORTED_VERSIONS = [...]string{"4", "5"}

type DHTMessage struct {
	Id        string "i"
	Query     string "q"
	Command   string "c"
	Arguments string "a"
	Payload   string "p"
}

type MSG_TYPE uint16

// Internal network packet type
const (
	MT_STRING     MSG_TYPE = 0 // String
	MT_INTRO               = 1 // Introduction packet
	MT_INTRO_REQ           = 2 // Request for introduction packet
	MT_NENC                = 3 // Not encrypted message
	MT_ENC                 = 4 // Encrypted message
	MT_PING                = 5 // Internal ping message for Proxies
	MT_XPEER_PING          = 6 // Crosspeer ping message
	MT_TEST                = 7 // Packet tests established connection
	MT_PROXY               = 8 // Information about proxy (forwarder)
	MT_BAD_TUN             = 9 // Notifies about dead tunnel
)

// List of commands used in DHT
const (
	CMD_CONN    string = "conn"
	CMD_FIND    string = "find"
	CMD_NODE    string = "node"
	CMD_PING    string = "ping"
	CMD_REGCP   string = "regcp"
	CMD_BADCP   string = "badcp"
	CMD_CP      string = "cp"
	CMD_NOTIFY  string = "notify"
	CMD_LOAD    string = "load"
	CMD_STOP    string = "stop"
	CMD_UNKNOWN string = "unk"
	CMD_DHCP    string = "dhcp"
	CMD_ERROR   string = "error"
)

const (
	DHT_ERROR_UNSUPPORTED string = "unsupported"
)

type (
	PeerState int
	PingType  uint16
)

// Peer state
const (
	P_INIT                  PeerState = iota // Peer has been added recently.
	P_REQUESTED_IP                    = iota // We know ID of a peer, but don't know it's IPs
	P_CONNECTING_DIRECTLY             = iota // Trying to establish a direct connection
	P_CONNECTED                       = iota // Connected, handshaked and operating normally
	P_HANDSHAKING                     = iota // Handshake requsted
	P_HANDSHAKING_FAILED              = iota // Handshake procedure failed
	P_WAITING_FORWARDER               = iota // Forwarder was requested
	P_HANDSHAKING_FORWARDER           = iota // Forwarder has been received and we're trying to handshake it
	P_DISCONNECT                      = iota // We're disconnecting
	P_STOP                            = iota // Peer has been stopped and now can be removed from list of peers
)

// Ping types
const (
	PING_REQ  PingType = 1
	PING_RESP PingType = 2
)

// Timeouts and retries
const (
	DHT_MAX_RETRIES         int           = 10
	DHCP_MAX_RETRIES        int           = 10
	PEER_PING_TIMEOUT       time.Duration = 5 * time.Second
	WAIT_PROXY_TIMEOUT      time.Duration = 2 * time.Second
	HANDSHAKE_PROXY_TIMEOUT time.Duration = 3 * time.Second
)
