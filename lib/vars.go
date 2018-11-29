package ptp

import (
	"time"
)

// PacketVersion is a version of packet used in DHT communication
const PacketVersion int32 = 20005

// SupportedVersion is a list of versions supported by DHT server
var SupportedVersion = [...]int32{20005, 200006}

// DHTBufferSize is a size of DHT buffer
const DHTBufferSize = 4096

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
	MsgTypeLatency           = 11 // Latency measurement
	MsgTypeComm              = 12 // Internal cross peer communication
)

// Common communication packet types
const (
	CommStatusReport uint16 = 0 // Status report between peers
	CommPing                = 1 // Ping packet
	CommLatency             = 2 // Latency packet
)

// IP communication packets
const (
	CommIPSubnet   uint16 = 10 // Request subnet information from peer
	CommIPInfo            = 11 // Ask peer if it knows specified IP
	CommIPSet             = 12 // Notify peer that this peer is now available over specified IP
	CommIPConflict        = 13 // Notify peer that his IP is in conflict
)

// Discovery communication packets
const (
	CommDiscoveryInit uint16 = 20 // Initiate connection with discovery service
	CommDiscoveryFind        = 21 // Find request
)

// Network Constants
const (
	MagicCookie uint16 = 0xabcd
	HeaderSize  int    = 10
)

// Network Variables

// LatencyProxyHeader used as a header of proxy request
var LatencyProxyHeader = []byte{0xfa, 0xca, 0x13, 0x15}

// LatencyRequestHeader used as a header when sending latency request
var LatencyRequestHeader = []byte{0xde, 0xad, 0xde, 0xda}

// LatencyResponseHeader used as a header when sending latency response
var LatencyResponseHeader = []byte{0xad, 0xde, 0xad, 0xde}

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
	DHTMaxRetries                  int           = 10
	DHCPMaxRetries                 int           = 10
	PeerPingTimeout                time.Duration = time.Second * 1
	WaitProxyTimeout               time.Duration = time.Second * 5
	HandshakeProxyTimeout          time.Duration = time.Second * 3
	EndpointPingInterval           time.Duration = time.Millisecond * 7000
	EndpointTimeout                time.Duration = time.Millisecond * 15000 // Must be greater than EndpointPingInterval
	ProxyLatencyRequestInterval    time.Duration = time.Second * 15         // How often we should update latency with proxies
	EndpointLatencyRequestInterval time.Duration = time.Second * 15         // How often we should update latency with endpoints
	UDPHolePunchTimeout            time.Duration = time.Millisecond * 20000 // How long we will for udp hole punching to finish
)
