package ptp

const PACKET_VERSION string = "3"

var SUPPORTED_VERSIONS = [...]string{"3"}

type DHTMessage struct {
	Id        string "i"
	Query     string "q"
	Command   string "c"
	Arguments string "a"
	Payload   string "p"
}

// TODO: Modify these structures
type DHTRequest struct {
	Id        string "i"
	Query     string "q"
	Command   string "c"
	Arguments string "a"
}

type DHTResponse struct {
	Id      string "i"
	Dest    string "h"
	Command string "c"
}

type MSG_TYPE uint16

// Internal network packet type
const (
	MT_STRING    MSG_TYPE = 0 + iota // String
	MT_INTRO              = 1        // Introduction packet
	MT_INTRO_REQ          = 2        // Request for introduction packet
	MT_NENC               = 3        // Not encrypted message
	MT_ENC                = 4        // Encrypted message
	MT_PING               = 5        // Internal ping message
	MT_TEST               = 6        // Packet tests established connection
	MT_PROXY              = 7        // Information about proxy (forwarder)
	MT_BAD_TUN            = 8        // Notifies about dead tunnel
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
)

// Peer state
const (
	P_INIT                  PeerState = 0 + iota
	P_CONNECTED                       = 1
	P_HANDSHAKING                     = 2
	P_HANDSHAKING_FAILED              = 3
	P_WAITING_FORWARDER               = 4
	P_HANDSHAKING_FORWARDER           = 5
	P_DISCONNECT                      = 6
)
