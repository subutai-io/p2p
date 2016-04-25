APP=p2p
CC=go
PACK=goupx
VERSION=$(shell git describe)
OS=$(shell uname -s)
ARCH=$(shell uname -m)

all: pack

$(APP): help.go instance.go main.go
	$(CC) build -ldflags="-w -s -X main.VERSION=$(VERSION)" -o $@ -v $^

pack: $(APP)
	$(PACK) $(APP)

clean:
	-rm -f $(APP)
	-rm -f $(APP)-*-v*
	-rm -f $(APP)-v*

test:  $(APP)
	go test ./...

release: $(APP)
release: pack
release:
	-cp $(APP) $(APP)-$(OS)-$(ARCH)-$(VERSION)
