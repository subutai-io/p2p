package ptp

import (
	"fmt"
	"net"
	"testing"
	"reflect"
	"sync"
)

func TestGenerateMac(t *testing.T) {
	macs := make(map[string]net.HardwareAddr)

	for i := 0; i < 10000; i++ {
		smac, mac := GenerateMAC()
		if smac == "" {
			t.Errorf("Failed to generate mac")
			return
		}
		_, e := macs[smac]
		if e {
			t.Errorf("Same MAC was generated")
			return
		}
		macs[smac] = mac
	}
}

func TestFindNetworkAddresses(t *testing.T) {
	ptp := new(PeerToPeer)
	ptp.FindNetworkAddresses()
	fmt.Printf("%+v\n", ptp.LocalIPs)
}
/*
Generated TestGenerateMAC
Generated TestGenerateToken
Generated Test_isPrivateIP
Generated TestStringifyState
Generated TestIsInterfaceLocal
Generated TestPeerToPeer_FindNetworkAddresses
Generated Test_min
package ptp

import (
	"fmt"
	"net"
	"reflect"
	"sync"
	"testing"
)
*/

/*
func TestGenerateMac(t *testing.T) {
	macs := make(map[string]net.HardwareAddr)

	for i := 0; i < 10000; i++ {
		smac, mac := GenerateMAC()
		if smac == "" {
			t.Errorf("Failed to generate mac")
			return
		}
		_, e := macs[smac]
		if e {
			t.Errorf("Same MAC was generated")
			return
		}
		macs[smac] = mac
	}
}

func TestFindNetworkAddresses(t *testing.T) {
	ptp := new(PeerToPeer)
	ptp.FindNetworkAddresses()
	fmt.Printf("%+v\n", ptp.LocalIPs)
}
*/

func TestGenerateMAC(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		want1 net.HardwareAddr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GenerateMAC()
			if got != tt.want {
				t.Errorf("GenerateMAC() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GenerateMAC() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateToken(); got != tt.want {
				t.Errorf("GenerateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPrivateIP(t *testing.T) {
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isPrivateIP(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("isPrivateIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isPrivateIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringifyState(t *testing.T) {
	type args struct {
		state PeerState
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringifyState(tt.args.state); got != tt.want {
				t.Errorf("StringifyState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInterfaceLocal(t *testing.T) {
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInterfaceLocal(tt.args.ip); got != tt.want {
				t.Errorf("IsInterfaceLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_FindNetworkAddresses(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			p.FindNetworkAddresses()
		})
	}
}

func Test_min(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}
