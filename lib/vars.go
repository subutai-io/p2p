package ptp

import (
	"time"
)

// PacketVersion is a version of packet used in DHT communication
const PacketVersion int32 = 20005

// SupportedVersion is a list of versions supported by DHT server
var SupportedVersion = [...]int32{20005, 200006}

// DHTBufferSize is a size of DHT buffer
const DHTBufferSize = 1024

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
	PeerStateInit             PeerState = 1  // Peer has been added recently.
	PeerStateRequestedIP                = 2  // We know ID of a peer, but don't know it's IPs
	PeerStateRequestingProxy            = 3  // Requesting proxies for this peer
	PeerStateWaitingForProxy            = 4  // Waiting for proxies
	PeerStateWaitingToConnect           = 5  // Waiting for other peer to start establishing connection
	PeerStateConnecting                 = 6  // Trying to establish connection
	PeerStateConnected                  = 7  // Connected, handshaked and operating normally
	PeerStateDisconnect                 = 8  // We're disconnecting
	PeerStateStop                       = 9  // Peer has been stopped and now can be removed from list of peers
	PeerStateCooldown                   = 10 // Peer is in cooldown mode
)

// Timeouts and retries
const (
	DHTMaxRetries         int           = 10
	DHCPMaxRetries        int           = 10
	PeerPingTimeout       time.Duration = time.Second * 1
	WaitProxyTimeout      time.Duration = time.Second * 5
	HandshakeProxyTimeout time.Duration = time.Second * 3
)
