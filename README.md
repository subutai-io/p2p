P2P Cloud
===================

P2P Cloud project allows users to build their private networks. 

Running
-------------------

p2p is managed by a daemon that controls every instance of your private networks (if you're participating in a different networks at the same time). To start a daemon simply run p2p with -daemon flag. Note, that application will run in a foreground mode. 

```
p2p -daemon
```

Now you can start manage the daemon with p2p command line interface. To start a new network or join existing you should run p2p application with a -start flag.

```
p2p -start -ip 10.10.10.1 -hash UNIQUE_STRING_IDENTIFIER
```

You should specify an IP address which will be used by your virtual network interface. All the participants should have an agreement on ranges of IP addresses they're using. In the future this will become unnecessary, because DHCP-like service will be implemented.

With a -hash flag user should specify a unique name of his network. 

Instance of P2P network can be stopped with use of -stop flag

```
p2p -stop -hash UNIQUE_STRING_IDENTIFIER
```

Development & Branching Model
-------------------

* 'master' is always stable. 
* 'dev' contains latest development snapshot that is under heavy testing

Dependencies
-------------------

* github.com/danderson/tuntap

