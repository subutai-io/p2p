CC=go
PACK=goupx
BUILD=$(shell git describe)
VERSION=$(shell cat VERSION)
OS=$(shell uname -s)
ARCH=$(shell uname -m)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
NAME_PREFIX=p2p
NAME_BASE=p2p
SOURCES=instance.go main.go rest.go start.go stop.go show.go set.go status.go debug.go daemon.go dht_connection.go dht_router.go
DOMAIN=subutai.io

sinclude config.make
ifdef DHT_ENDPOINTS
	DHT=$(DHT_ENDPOINTS)
else
	DHT=dht
endif
APP=$(NAME_PREFIX)

build: directories
build: bin/$(APP)
linux: bin/$(APP)
windows: bin/$(APP).exe
macos: bin/$(APP)_osx
all: linux windows macos

bin/$(APP): $(SOURCES) service_posix.go
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	GOOS=linux $(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)$(BRANCH_POSTFIX) -X main.TargetURL=$(DHT) -X main.BuildID=$(BUILD) -X main.DefaultLog=$(LOG_LEVEL)" -o $@ -v $^

bin/$(APP).exe: $(SOURCES) service_windows.go
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	GOOS=windows $(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)$(BRANCH_POSTFIX) -X main.TargetURL=$(DHT) -X main.BuildID=$(BUILD) -X main.DefaultLog=$(LOG_LEVEL)" -o $@ -v $^
	
bin/$(APP)_osx: $(SOURCES) service_posix.go
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	GOOS=darwin $(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)$(BRANCH_POSTFIX) -X main.TargetURL=$(DHT) -X main.BuildID=$(BUILD) -X main.DefaultLog=$(LOG_LEVEL)" -o $@ -v $^

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
	go test -v github.com/subutai-io/p2p/lib
	go test --bench . ./...

coverage:
	go test -coverprofile=coverage.txt -covermode=atomic github.com/subutai-io/p2p/lib

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

proto:
	protoc --go_out=import_path=protocol:. protocol/dht.proto