#!/bin/bash

BINARY=$1
KEY=$2

SSH_USER="ubuntu"
APP_NAME="p2p-cp"

if [ -z $BINARY ]; then
    echo "Specify path to a Control Peer binary"
    exit
fi

if [ -z $KEY ]; then
    echo "Specify path to PEM Keyfile"
    exit
fi

if [ ! -f $BINARY ]; then
    echo "Failed to find Control Peer binary at specified path: $BINARY"
    exit
fi

if [ ! -f $KEY ]; then
    echo "$KEY PEM Key file not found"
    exit
fi

for host in 1 2 3 4 5; do
    ssh -i $KEY $SSH_USER@dht$host.subut.ai "killall -9 $APP_NAME"
    ssh -i $KEY $SSH_USER@dht$host.subut.ai "killall -9 cp"
    scp -i $KEY $BINARY $SSH_USER@dht$host.subut.ai:~
    ssh -n -f -i $KEY $SSH_USER@dht$host.subut.ai "~/$APP_NAME 2> /home/ubuntu/p2p-cp.log &" &
done
