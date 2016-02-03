APP=p2p
CP_APP=cp
CC=go
PACK=goupx
VERSION=$(shell git describe)

all: pack

$(APP): p2p.go packet.go main.go
	$(CC) build -ldflags="-w -s" -o $@ -v $^

$(CP_APP): p2p-cp/cp.go p2p-cp/proxy.go
	$(CC) build -ldflags="-w -s" -o $@ -v $^

pack: $(APP) $(CP_APP)
	$(PACK) $(APP)
	$(PACK) $(CP_APP)

clean:
	-rm -f $(APP)
	-rm -f $(CP_APP)
	-rm -f $(APP)-v*
	-rm -f $(CP_APP)-v*

test:  $(APP) $(CP_APP)
	go test ./...

release: $(APP) $(CP_APP)
release: pack
release:
	-cp $(APP) $(APP)-$(VERSION)
	-cp $(CP_APP) $(CP_APP)-$(VERSION)
