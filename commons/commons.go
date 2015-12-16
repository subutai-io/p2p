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
	MT_STRING MSG_TYPE = 0 // String
	MT_INTRO  MSG_TYPE = 1 // Introduction packet
	//todo add types
)
