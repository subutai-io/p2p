CC=go
PACK=goupx
BUILD=$(shell git describe)
VERSION=$(shell cat VERSION)
OS=$(shell uname -s)
ARCH=$(shell uname -m)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
NAME_PREFIX=p2p
NAME_BASE=p2p
SOURCES=help.go instance.go main.go rest.go start.go stop.go show.go set.go status.go debug.go
DHT=mdht.subut.ai:6881
ifeq ($(BRANCH),HEAD)
	DHT=mdht.subut.ai:6881
	SCHEME=
else
	SCHEME=-$(BRANCH)
endif
ifeq ($(BRANCH),dev)
	DHT=18.195.169.215:6881
endif
ifeq ($(BRANCH),master)
	DHT=54.93.172.70:6881
endif
ifeq ($(BRANCH),sysnet)
	DHT=18.195.169.215:6881
endif
sinclude config.make
APP=$(NAME_BASE)


build: $(APP)
ifdef UPX_BIN
release: pack
endif

linux: $(APP)
windows: $(APP).exe
macos: $(APP)_osx


all: linux windows macos

$(APP): $(SOURCES) service_posix.go
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	$(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)$(SCHEME) -X main.DefaultDHT=$(DHT) -X main.BuildID=$(BUILD)" -o $@ -v $^

$(APP).exe: $(SOURCES) service_windows.go
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	GOOS=windows $(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)$(SCHEME) -X main.DefaultDHT=$(DHT) -X main.BuildID=$(BUILD)" -o $@ -v $^
	
$(APP)_osx: $(SOURCES) service_posix.go
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	GOOS=darwin $(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)$(SCHEME) -X main.DefaultDHT=$(DHT) -X main.BuildID=$(BUILD)" -o $@ -v $^

ifdef UPX_BIN
pack: $(APP)
	$(PACK) $(APP)
endif

clean:
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
	-rm -f config.make

ifeq ($(BUILD_DEB), 1)
test:  $(APP)
	go test ./...
else
test: skip-test
endif

skip-test: $(APP)
	@echo "Test skipped"

release: build
release:
	-mv $(APP) $(APP)-$(OS)-$(ARCH)-$(VERSION)

install: 
	@mkdir -p $(DESTDIR)/opt/subutai/bin
	@cp $(APP) $(DESTDIR)/opt/subutai/bin/$(NAME_PREFIX)

uninstall:
	@rm -f $(DESTDIR)/bin/$(NAME_PREFIX)

ifeq ($(BUILD_DEB), 1)
debian: *.deb

*.deb:
	debuild --preserve-env -B -d

debian-source: *.changes

*.changes:
	debuild --preserve-env -S -d
endif

snapcraft: $(SOURCES) service_posix.go
	$(eval export GOPATH=$(shell pwd)/../go)
	$(eval export GOBIN=$(shell pwd)/../go/bin)
	@if [ ! -d "$(GOPATH)/src/github.com/subutai-io/p2p" ]; then mkdir -p $(GOPATH)/src/github.com/subutai-io/; ln -s $(shell pwd) $(GOPATH)/src/github.com/subutai-io/p2p; fi
	$(CC) get -d
	$(CC) build -ldflags="-r /apps/subutai/current/lib -w -s -X main.AppVersion=$(VERSION)$(SCHEME) -X main.DefaultDHT=$(DHT) -X main.BuildID=$(BUILD)" -o $(APP) -v $^
