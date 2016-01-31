APP=p2p
CP_APP=cp
CC=go
PACK=goupx

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

test:  $(APP) $(CP_APP)
	go test ./...
