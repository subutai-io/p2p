package ptp

// TODO: Modify these structures
type DHTRequest struct {
	Id      string "i"
	Hash    string "h"
	Command string "c"
	Port    string "p"
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
	MT_INTRO_REQ          = 2
	MT_NENC               = 3 // Not encrypted message
	MT_ENC                = 4 // Encrypted message
	MT_PING               = 5 // Internal ping message
	MT_TEST               = 6
	MT_PROXY              = 7
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
)
