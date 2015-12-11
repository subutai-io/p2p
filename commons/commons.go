package commons

type DHTRequest struct {
	Id      string "i"
	Hash    string "h"
	Command string "c"
}

type DHTResponse struct {
	Id      string "i"
	Dest    string "h"
	Command string "c"
}
