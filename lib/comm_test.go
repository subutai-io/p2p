package ptp

import (
	"encoding/binary"
	"net"
	"reflect"
	"testing"
)

func Test_commPacketCheck(t *testing.T) {
	type args struct {
		data []byte
	}

	ut := "123e4567-e89b-12d3-a456-426655440000"

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"nil case", args{}, true},
		{"small size", args{[]byte{0x01}}, true},
		{"passing", args{[]byte(ut)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := commPacketCheck(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("commPacketCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_commStatusReportHandler(t *testing.T) {
	type args struct {
		data []byte
		p    *PeerToPeer
	}

	ut := "123e4567-e89b-12d3-a456-426655440000"

	ptp := new(PeerToPeer)

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"nil case", args{nil, ptp}, nil, true},
		{"small size", args{[]byte{0x01}, ptp}, nil, true},
		{"passing", args{[]byte(ut), ptp}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := commStatusReportHandler(tt.args.data, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("commStatusReportHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commStatusReportHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commSubnetInfoHandler(t *testing.T) {
	type args struct {
		data []byte
		p    *PeerToPeer
	}

	ut := "123e4567-e89b-12d3-a456-426655440000"

	ptp := new(PeerToPeer)

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"nil case", args{nil, ptp}, nil, true},
		{"small size", args{[]byte{0x01}, ptp}, nil, true},
		{"passing", args{[]byte(ut), ptp}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := commSubnetInfoHandler(tt.args.data, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("commSubnetInfoHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commSubnetInfoHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commIPInfoHandler(t *testing.T) {
	type args struct {
		data []byte
		p    *PeerToPeer
	}

	ip := net.ParseIP("127.0.0.1")

	ut := "123e4567-e89b-12d3-a456-426655440000"

	ptp0 := new(PeerToPeer)

	ptp1 := new(PeerToPeer)
	ptp1.Peers = new(PeerList)
	ptp1.Peers.Init()

	ptp2 := new(PeerToPeer)
	ptp2.Peers = new(PeerList)
	ptp2.Peers.Init()
	ptp2.Interface, _ = newTAP("ip", "10.10.10.1", "00:11:22:33:44:55", "255.255.255.0", 1500, false)

	ptp3 := new(PeerToPeer)
	ptp3.Peers = new(PeerList)
	ptp3.Peers.Init()
	ptp3.Peers.peers["127.0.0.1"] = &NetworkPeer{
		PeerLocalIP: ip,
	}
	ptp3.Interface, _ = newTAP("ip", "10.10.10.1", "00:11:22:33:44:55", "255.255.255.0", 1500, false)

	d0 := append([]byte{0x00, 0x01}, []byte(ut)...)

	d1 := make([]byte, 42)
	copy(d1[0:36], ut)
	copy(d1[36:40], ip.To4())
	binary.BigEndian.PutUint16(d1[40:42], uint16(0))

	d2 := make([]byte, 42)
	copy(d2[0:36], ut)
	copy(d2[36:40], ip.To4())
	binary.BigEndian.PutUint16(d2[40:42], uint16(1))

	d3 := make([]byte, 41)
	copy(d3[0:36], ut)
	copy(d3[36:40], ip.To4())
	d3[40] = 0x0f

	d4 := make([]byte, 40)
	copy(d4[0:36], ut)
	copy(d4[36:40], ip.To4())

	result := ut + string(ip.To4())

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"nil case", args{}, nil, true},
		{"nil ptp", args{d0, nil}, nil, true},
		{"nil peer list", args{d0, ptp0}, nil, true},
		{"nil interface", args{d0, ptp1}, nil, true},
		{"too small", args{d0, ptp2}, nil, true},
		{"42 size>0", args{d1, ptp2}, nil, false},
		{"42 size>1", args{d2, ptp2}, nil, false},
		{"41 size", args{d3, ptp2}, nil, true},
		{"40 size", args{d4, ptp3}, append([]byte(result), []byte{0x00, 0x01}...), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := commIPInfoHandler(tt.args.data, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("commIPInfoHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commIPInfoHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commIPSetHandler(t *testing.T) {
	type args struct {
		data []byte
		p    *PeerToPeer
	}

	ut := "123e4567-e89b-12d3-a456-426655440000"

	ptp := new(PeerToPeer)

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"nil case", args{nil, ptp}, nil, true},
		{"small size", args{[]byte{0x01}, ptp}, nil, true},
		{"passing", args{[]byte(ut), ptp}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := commIPSetHandler(tt.args.data, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("commIPSetHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commIPSetHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commIPConflictHandler(t *testing.T) {
	type args struct {
		data []byte
		p    *PeerToPeer
	}

	ut := "123e4567-e89b-12d3-a456-426655440000"

	ptp := new(PeerToPeer)

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"nil case", args{nil, ptp}, nil, true},
		{"small size", args{[]byte{0x01}, ptp}, nil, true},
		{"passing", args{[]byte(ut), ptp}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := commIPConflictHandler(tt.args.data, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("commIPConflictHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commIPConflictHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
