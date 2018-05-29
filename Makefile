CC=go
PACK=goupx
BUILD=$(shell git describe)
VERSION=$(shell cat VERSION)
OS=$(shell uname -s)
ARCH=$(shell uname -m)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
NAME_PREFIX=p2p
NAME_BASE=p2p
APP_DIR=app
SOURCES=$(APP_DIR)/instance.go \
		$(APP_DIR)/main.go \
		$(APP_DIR)/rest.go \
		$(APP_DIR)/start.go \
		$(APP_DIR)/stop.go \
		$(APP_DIR)/show.go \
		$(APP_DIR)/set.go \
		$(APP_DIR)/status.go \
		$(APP_DIR)/debug.go \
		$(APP_DIR)/daemon.go \
		$(APP_DIR)/dht_connection.go \
		$(APP_DIR)/dht_router.go

sinclude config.make
ifdef DHT_ENDPOINTS
	DHT=$(DHT_ENDPOINTS)
else
	DHT=eu0.cdn.subutai.io:6881
endif
APP=$(NAME_PREFIX)

SNAPDHT=eu0.cdn.subutai.io:6881
ifeq ($(BRANCH),dev)
	SNAPDHT=eu0.devcdn.subutai.io:6881
endif
ifeq ($(BRANCH),master)
	SNAPDHT=eu0.mastercdn.subutai.io:6881
endif
ifeq ($(BRANCH),sysnet)
	SNAPDHT=eu0.sysnetcdn.subutai.io:6881
endif

build: directories
build: bin/$(APP)
linux: bin/$(APP)
windows: bin/$(APP).exe
macos: bin/$(APP)_osx
all: linux windows macos

bin/$(APP): $(SOURCES) $(APP_DIR)/service_posix.go
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	GOOS=linux $(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)$(BRANCH_POSTFIX) -X main.DefaultDHT=$(DHT) -X main.BuildID=$(BUILD) -X main.DefaultLog=$(LOG_LEVEL)" -o $@ -v $^

bin/$(APP).exe: $(SOURCES) $(APP_DIR)/service_windows.go
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	GOOS=windows $(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)$(BRANCH_POSTFIX) -X main.DefaultDHT=$(DHT) -X main.BuildID=$(BUILD) -X main.DefaultLog=$(LOG_LEVEL)" -o $@ -v $^
	
bin/$(APP)_osx: $(SOURCES) $(APP_DIR)/service_posix.go
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	GOOS=darwin $(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)$(BRANCH_POSTFIX) -X main.DefaultDHT=$(DHT) -X main.BuildID=$(BUILD) -X main.DefaultLog=$(LOG_LEVEL)" -o $@ -v $^

clean:
	-rm -f bin/$(APP)
	-rm -f bin/$(APP).exe
	-rm -f bin/$(APP)_osx
	-rm -f $(APP)
	-rm -f $(APP)_osx
	-rm -f $(APP).exe
	-rm -f $(APP)-$(OS)*
	-rm -f $(NAME_PREFIX)
	-rm -f $(NAME_PREFIX)_osx
	-rm -f $(NAME_PREFIX).exe
	-rm -f $(NAME_PREFIX)-$(OS)*
	-rm -rf debian/extra-code/*

mrproper: clean
mrproper:
	-rm -rf bin
	-rm -f config.make

test:
	go test -v github.com/subutai-io/p2p
	go test --bench . ./...

coverage:
	go test -coverprofile=coverage.txt -covermode=atomic 

release: build
release:
	-mv $(APP) $(APP)-$(OS)-$(ARCH)-$(VERSION)

install: 
	@mkdir -p $(DESTDIR)/opt/subutai/bin
	@cp $(APP) $(DESTDIR)/opt/subutai/bin/$(NAME_PREFIX)

uninstall:
	@rm -f $(DESTDIR)/bin/$(NAME_PREFIX)

directories:
	@mkdir -p bin

snapcraft: $(SOURCES) service_posix.go
	$(eval export GOPATH=$(shell pwd)/../go)
	$(eval export GOBIN=$(shell pwd)/../go/bin)
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	$(CC) get -d
	$(CC) get -u github.com/golang/protobuf/proto
	$(CC) build -ldflags="-r /apps/subutai/current/lib -w -s -X main.AppVersion=$(VERSION)$(BRANCH_POSTFIX) -X main.DefaultDHT=$(SNAPDHT) -X main.BuildID=$(BUILD) -X main.DefaultLog=$(LOG_LEVEL)" -o $(APP) -v $^

proto:
	protoc --go_out=import_path=protocol:. protocol/dht.proto
