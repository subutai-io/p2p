#!/bin/bash
# This script read file provided in argument with list of username@host
# each one on a new line. Then it uploads p2p application to each of them and runs it with

SSH_USER="root"
BINARY=$1
CONFIG=$2
APP_NAME="p2p"

if [ -z $BINARY ]; then
    echo "Specify path to a binary"
    exit
fi

if [ ! -f $BINARY ]; then
    echo "Failed to find binary at specified path: $BINARY"
    exit
fi

if [ -z $CONFIG ]; then
    echo "Specify path to a config YAML"
    exit
fi

if [ ! -f $CONFIG ]; then
    echo "Failed to find config YAML: $CONFIG"
    exit
fi

for HOST_ID in 1 2 3 4; do
    ssh $SSH_USER@bigdata$HOST_ID "killall -9 $APP_NAME"
    scp $BINARY $SSH_USER@bigdata$HOST_ID:~
    scp $CONFIG $SSH_USER@bigdata$HOST_ID:~
    ssh $SSH_USER@bigdata$HOST_ID "~/$APP_NAME -ip 10.10.10.10$HOST_ID -mask 255.255.255.0 -dev tun0 -hash somerandomhash 2> /var/log/$APP_NAME.log" &
done
