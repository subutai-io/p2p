#!/bin/bash

cd ../
go build
scp p2p user@192.168.56.102:~/
scp p2p user@192.168.56.103:~/

cd p2p-cp
go build
scp p2p-cp user@192.168.56.102:~/
scp p2p-cp user@192.168.56.103:~/
