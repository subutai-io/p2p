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
