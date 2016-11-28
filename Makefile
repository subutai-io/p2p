CC=go
PACK=goupx
VERSION=$(shell git describe)
OS=$(shell uname -s)
ARCH=$(shell uname -m)
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

$(APP): help.go instance.go main.go
	$(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)" -o $@ -v $^

$(APP).exe:
	GOOS=windows $(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)" -o $@ -v $^
	
$(APP)_osx:
	GOOS=darwin $(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)" -o $@ -v $^

pack: $(APP)
	$(PACK) $(APP)

clean:
	-rm -f $(APP)
	-rm -f $(APP)_osx
	-rm -f $(APP).exe
	-rm -f $(APP)-$(OS)*
	-rm -f $(NAME_PREFIX)
	-rm -f $(NAME_PREFIX)_osx
	-rm -f $(NAME_PREFIX).exe
	-rm -f $(NAME_PREFIX)-$(OS)*

mrproper: clean
mrproper:
	-rm -f config.make

test:  $(APP)
	go test ./...

release: build
release:
	-mv $(APP) $(APP)-$(OS)-$(ARCH)-$(VERSION)

install: 
	@mkdir -p $(DESTDIR)/bin
	@cp $(APP) $(DESTDIR)/bin

ifdef BUILD_DEBIAN
debian:
	debuild -B -d
endif
