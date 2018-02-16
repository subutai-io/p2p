package ptp

import (
	"time"
)

// PacketVersion is a version of packet used in DHT communication
const PacketVersion int32 = 20005

// SupportedVersion is a list of versions supported by DHT server
var SupportedVersion = [...]int32{20005, 200006}

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
	DhtCmdState       string = "state"
)

const (
	// DhtErrorUnsupported - Unsupported version
	DhtErrorUnsupported string = "unsupported"
)

// PeerState - current state of the peer
type PeerState int

// Peer state
const (
	PeerStateInit             PeerState = iota // Peer has been added recently.
	PeerStateRequestedIP                = iota // We know ID of a peer, but don't know it's IPs
	PeerStateRequestingProxy            = iota // Requesting proxies for this peer
	PeerStateWaitingForProxy            = iota // Waiting for proxies
	PeerStateWaitingToConnect           = iota // Waiting for other peer to start establishing connection
	PeerStateConnecting                 = iota // Trying to establish connection
	PeerStateRouting                    = iota // (Re)Routing
	PeerStateConnected                  = iota // Connected, handshaked and operating normally
	PeerStateDisconnect                 = iota // We're disconnecting
	PeerStateStop                       = iota // Peer has been stopped and now can be removed from list of peers
	PeerStateCooldown                   = iota
)

// Timeouts and retries
const (
	DHTMaxRetries         int           = 10
	DHCPMaxRetries        int           = 10
	PeerPingTimeout       time.Duration = time.Second * 1
	WaitProxyTimeout      time.Duration = time.Second * 5
	HandshakeProxyTimeout time.Duration = time.Second * 3
)
