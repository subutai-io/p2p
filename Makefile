CC=go
PACK=goupx
BUILD=$(shell git describe)
VERSION=$(shell cat VERSION)
OS=$(shell uname -s)
ARCH=$(shell uname -m)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
NAME_PREFIX=p2p
NAME_BASE=p2p
SOURCES=help.go instance.go main.go rest.go start.go stop.go show.go set.go status.go debug.go daemon.go
#DHT=mdht.subut.ai:6881
#ifeq ($(BRANCH),HEAD)
#	DHT=mdht.subut.ai:6881
#	SCHEME=
#else
#	SCHEME=-$(BRANCH)
#endif


sinclude config.make
ifdef DHT_ENDPOINTS
	DHT=$(DHT_ENDPOINTS)
else
	DHT=mdht.subut.ai:6881
endif
APP=$(NAME_PREFIX)

SNAPDHT=mdht.subut.ai:6881
ifeq ($(BRANCH),dev)
	SNAPDHT=18.195.169.215:6881
endif
ifeq ($(BRANCH),master)
	SNAPDHT=54.93.172.70:6881
endif
ifeq ($(BRANCH),sysnet)
	SNAPDHT=18.195.169.215:6881
endif

build: directories
build: bin/$(APP)
linux: bin/$(APP)
windows: bin/$(APP).exe
macos: bin/$(APP)_osx
all: linux windows macos

bin/$(APP): $(SOURCES) service_posix.go
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	$(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)$(BRANCH_POSTFIX) -X main.DefaultDHT=$(DHT) -X main.BuildID=$(BUILD)" -o $@ -v $^

bin/$(APP).exe: $(SOURCES) service_windows.go
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	GOOS=windows $(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)$(BRANCH_POSTFIX) -X main.DefaultDHT=$(DHT) -X main.BuildID=$(BUILD)" -o $@ -v $^
	
bin/$(APP)_osx: $(SOURCES) service_posix.go
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	GOOS=darwin $(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)$(BRANCH_POSTFIX) -X main.DefaultDHT=$(DHT) -X main.BuildID=$(BUILD)" -o $@ -v $^

clean:
	-rm -f bin/$(APP)
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

test:  $(APP)
	go test -v ./...

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
	$(CC) build -ldflags="-r /apps/subutai/current/lib -w -s -X main.AppVersion=$(VERSION)$(BRANCH_POSTFIX) -X main.DefaultDHT=$(SNAPDHT) -X main.BuildID=$(BUILD)" -o $(APP) -v $^
