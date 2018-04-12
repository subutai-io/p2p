// +build gofuzz

package ptp

func Fuzz(data []byte) int {
	arp := new(ARPPacket)
	err := arp.UnmarshalARP(data)
	if err != nil {
		return 0
	}
	return 1
}
