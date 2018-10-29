package main

import (
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestRestore_init(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}
	type args struct {
		filepath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"No save file", fields{}, args{filepath: ""}, false},
		{"Normal save file", fields{}, args{filepath: "/tmp/t1"}, false},
		{"Bad filename", fields{}, args{filepath: "/"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			if err := r.init(tt.args.filepath); (err != nil) != tt.wantErr {
				t.Errorf("Restore.init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestore_save(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Bad filename", fields{filepath: "/"}, true},
		{"Normal filename", fields{filepath: "/tmp/t2"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			if err := r.save(); (err != nil) != tt.wantErr {
				t.Errorf("Restore.save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestore_load(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}

	presaved1 := "- empty"
	f1, _ := os.OpenFile("/tmp/restore-load-test1", os.O_CREATE|os.O_RDWR, 0700)
	f1.Write([]byte(presaved1))
	f1.Close()

	f2, _ := os.OpenFile("/tmp/restore-load-test2", os.O_CREATE|os.O_RDWR, 0700)
	f2.Write([]byte(""))
	f2.Close()

	f3, _ := os.OpenFile("/tmp/restore-load-test3", os.O_CREATE|os.O_RDWR, 0700)
	f3.Write([]byte("10.132.190.1~~p2p8~swarm-5565bd8c-d570-4fc8-bf39-eca47b918585~~~630a05b329266ac8a69adfe7d2d87003~1540794362~0~0|||10.243.26.1~~p2p4~swarm-9de835d3-eae8-49c3-9d61-1e65767618d9~~~3ea66ea049506bc018d33f82ec3b3cce~1527854837~0~0|||10.54.31.1~~p2p2~swarm-f98003ee-27c4-41ef-bc29-b6a1da5d779b~~~998e37ba12867e1dc6d074f0d3313ea9~1540710496~0~0|||10.36.162.2~~p2p6~swarm-ef0e47f8-6821-44e0-b821-ce909d661c4e~~~be2833276389bc76528cb8b227b50db5~1540750940~0~0|||10.205.196.1~~p2p9~swarm-2777bdf0-ede6-4c0d-b352-a3c474e7d7e0~~~48bd3833d2fad84e374b98c2e7c9f6ef~1539596713~0~0|||10.159.190.1~~p2p7~swarm-2120b53e-cd49-4eef-8927-24273686f3ab~~~0c46bd9eec679d4b97d26cb8a90d15b3~1539600413~0~0"))
	f3.Close()

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Bad filename", fields{filepath: "/"}, true},
		{"Normal file", fields{filepath: "/tmp/restore-load-test1"}, false},
		{"Empty file", fields{filepath: "/tmp/restore-load-test2"}, false},
		{"Old format", fields{filepath: "/tmp/restore-load-test3"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			if err := r.load(); (err != nil) != tt.wantErr {
				t.Errorf("Restore.load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestore_addInstance(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}
	type args struct {
		inst *P2PInstance
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Simple instance", fields{}, args{&P2PInstance{ID: "mh"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			if err := r.addInstance(tt.args.inst); (err != nil) != tt.wantErr {
				t.Errorf("Restore.addInstance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestore_addEntry(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}
	type args struct {
		entry saveEntry
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"New entry", fields{}, args{saveEntry{Hash: "hash"}}, false},
		{"Existing entry", fields{entries: []saveEntry{{Hash: "hash"}}}, args{saveEntry{Hash: "hash"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			if err := r.addEntry(tt.args.entry); (err != nil) != tt.wantErr {
				t.Errorf("Restore.addEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestore_removeEntry(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Non-existing entry", fields{}, args{"hash"}, true},
		{"Existing entry", fields{entries: []saveEntry{{Hash: "hash"}}}, args{"hash"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			if err := r.removeEntry(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Restore.removeEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestore_bumpInstance(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Non-existing entry", fields{}, args{"hash"}, true},
		{"Existing entry", fields{entries: []saveEntry{{Hash: "hash"}}}, args{"hash"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			if err := r.bumpInstance(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Restore.bumpInstance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestore_disableStaleInstances(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}
	type args struct {
		inst *P2PInstance
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Pass", fields{}, args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			if err := r.disableStaleInstances(tt.args.inst); (err != nil) != tt.wantErr {
				t.Errorf("Restore.disableStaleInstances() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestore_encode(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}

	ee1 := `- ip: ""
  mac: ""
  dev: ""
  hash: hash
  keyfile: ""
  key: ""
  ttl: ""
  last_success: ""
  enabled: true
`

	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{"No entries", fields{}, nil, false},
		{"Single Disabled Entry", fields{entries: []saveEntry{{Hash: "hash"}}}, []byte{91, 93, 10}, false},
		{"Single Enabled Entry", fields{entries: []saveEntry{{Hash: "hash", Enabled: true}}}, []byte(ee1), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			got, err := r.encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Restore.encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Restore.encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestore_decode(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Failed case", fields{}, args{[]byte("/\\")}, true},
		{"Passing case", fields{}, args{[]byte("-")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			if err := r.decode(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Restore.decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestore_decodeInstances(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			if err := r.decodeInstances(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Restore.decodeInstances() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestore_get(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []saveEntry
	}{
		{"Passing", fields{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			if got := r.get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Restore.get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestore_isActive(t *testing.T) {
	type fields struct {
		entries  []saveEntry
		filepath string
		lock     sync.RWMutex
		active   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Active", fields{active: true}, true},
		{"Non active", fields{active: false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Restore{
				entries:  tt.fields.entries,
				filepath: tt.fields.filepath,
				lock:     tt.fields.lock,
				active:   tt.fields.active,
			}
			if got := r.isActive(); got != tt.want {
				t.Errorf("Restore.isActive() = %v, want %v", got, tt.want)
			}
		})
	}
}
