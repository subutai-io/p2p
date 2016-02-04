P2P Control Peer
===================

Control peer used by p2p client in a two ways: as a DHT node and as a proxy server (traffic forwarder). 

Distributed Hash Table
-------------------

Control Peer is not a Torrent-like DHT service, but it uses the same techniques to bring peers together. To start a new DHT node you need to run *cp* application with *-dht* argument and specify a UDP port which DHT will listen to.

```
cp -dht 6881
```

Traffic forwarding
-------------------

To forward traffic you need to run *cp* without any arguments. It will connect to a default DHT routers and register itself as a traffic forwarder that can be used by peers that is experiencing problems with cross-peer communication.
