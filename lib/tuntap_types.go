package ptp

const (
	flagTruncated = 0x1

	iffTun      = 0x1
	iffTap      = 0x2
	iffOneQueue = 0x2000
	iffnopi     = 0x1000
)

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}
