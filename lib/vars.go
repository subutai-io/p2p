package ptp

import (
	"time"
)

// PacketVersion is a version of packet used in DHT communication
const PacketVersion string = "5"

// SupportedVersion is a list of versions supported by DHT server
var SupportedVersion = [...]string{"6", "5"}

// DHTMessage is an unmarshaled DHT packet
type DHTMessage struct {
	ID        string "i"
	Query     string "q"
	Command   string "c"
	Arguments string "a"
	Payload   string "p"
}

// MsgType is a type of the message
type MsgType uint16

// Internal network packet type
const (
	MsgTypeString    MsgType = 0  // String
	MsgTypeIntro             = 1  // Introduction packet
	MsgTypeIntroReq          = 2  // Request for introduction packet
	MsgTypeNenc              = 3  // Not encrypted message
	MsgTypeEnc               = 4  // Encrypted message
	MsgTypePing              = 5  // Internal ping message for Proxies
	MsgTypeXpeerPing         = 6  // Crosspeer ping message
	MsgTypeTest              = 7  // Packet tests established connection
	MsgTypeProxy             = 8  // Information about proxy (forwarder)
	MsgTypeBadTun            = 9  // Notifies about dead tunnel
	MsgTypeConf              = 10 // Confirmation
)

// List of commands used in DHT
const (
	DhtCmdConn        string = "conn"
	DhtCmdFrwd        string = "frwd"
	DhtCmdFind        string = "find"
	DhtCmdNode        string = "node"
	DhtCmdPing        string = "ping"
	DhtCmdRegProxy    string = "regcp"
	DhtCmdBadProxy    string = "badcp"
	DhtCmdProxy       string = "cp"
	DhtCmdNotify      string = "notify"
	DhtCmdLoad        string = "load"
	DhtCmdStop        string = "stop"
	DhtCmdUnknown     string = "unk"
	DhtCmdDhcp        string = "dhcp"
	DhtCmdError       string = "error"
	DhtCmdUnsupported string = "unsupported"
)

const (
	// DhtErrorUnsupported - Unsupported version
	DhtErrorUnsupported string = "unsupported"
)

type (
	// PeerState - current state of the peer
	PeerState int
	// PingType - whether this is a request or a response
	PingType uint16
)

// Peer state
const (
	PeerStateInit                 PeerState = iota // Peer has been added recently.
	PeerStateRequestedIP                    = iota // We know ID of a peer, but don't know it's IPs
	PeerStateConnectingDirectly             = iota // Trying to establish a direct connection
	PeerStateConnected                      = iota // Connected, handshaked and operating normally
	PeerStateHandshaking                    = iota // Handshake requsted
	PeerStateHandshakingFailed              = iota // Handshake procedure failed
	PeerStateWaitingForwarder               = iota // Forwarder was requested
	PeerStateHandshakingForwarder           = iota // Forwarder has been received and we're trying to handshake it
	PeerStateDisconnect                     = iota // We're disconnecting
	PeerStateStop                           = iota // Peer has been stopped and now can be removed from list of peers
)

// Ping types
const (
	PingReq  PingType = 1
	PingResp PingType = 2
)

// Timeouts and retries
const (
	DHTMaxRetries         int           = 10
	DHCPMaxRetries        int           = 10
	PeerPingTimeout       time.Duration = time.Second * 15
	WaitProxyTimeout      time.Duration = time.Second * 5
	HandshakeProxyTimeout time.Duration = time.Second * 3
)
