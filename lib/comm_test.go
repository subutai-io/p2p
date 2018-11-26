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
	if ip == nil {
		t.FailNow()
	}

	ut := "123e4567-e89b-12d3-a456-426655440000"

	ptp0 := new(PeerToPeer)

	ptp1 := new(PeerToPeer)
	ptp1.Peers = new(PeerList)
	ptp1.Peers.Init()

	d0 := append([]byte{0x00, 0x01}, []byte(ut)...)

	d1 := make([]byte, 42)
	d1 = append([]byte(ut), ip.To4()...)
	binary.BigEndian.PutUint16(d1[40:42], uint16(1))

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"nil case", args{}, nil, true},
		{"nil ptp", args{d0, nil}, nil, true},
		{"nil peer list", args{d0, ptp0}, nil, true},
		{"too small", args{d0, ptp1}, nil, true},
		{"42 size", args{d1, ptp1}, nil, true},
		{"passing data", args{d1, ptp1}, nil, true},
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
