package commons

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

const (
	MT_STRING MSG_TYPE = 0 + iota // String
	MT_INTRO           = 1        // Introduction packet
	MT_NENC            = 2        // Not encrypted message
	MT_ENC             = 3        // Encrypted message
	MT_PING            = 4        // Internal ping message
	//todo add types
)
