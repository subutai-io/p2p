.. Subutai P2P documentation master file, created by
   sphinx-quickstart on Thu Aug 25 22:03:03 2016.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

Troubleshooting
=======================================

Collecting debug information
==========



**p2p debug** command helps to understand your setup and it's better to include output of this command in an issue.

Here is an example output of this command::

    DEBUG INFO:
    Number of gouroutines: 17
    Instances information:
    Hash: INTERNAL_DEV_TEST_SWARM_0747-1
    ID: 46bbb227-4d63-11e6-8718-022d0fae0f03
    Interface vptp1, HW Addr: 06:05:bb:2b:a6:42, IP: 10.10.1.1
    Peers:
        --- 4742e997-4d63-11e6-8e39-022d0fae0f03 ---
            HWAddr: 06:84:e6:ce:48:76
            IP: 10.10.1.245
            Endpoint: 52.28.78.136:53384
            Peer Address: 158.181.222.46:39400
            Proxy ID: 52
        --- End of 4742e997-4d63-11e6-8e39-022d0fae0f03 ---

There is a list of connected peers under "Peers" section.

* **IP** is an internal IP address
* **Endpoint** is an IP address used to communicate with this peer. If Endpoint equals Peer Address - that means peers are connected directly in LAN or over Internet. If it's not, that means that traffic forwarder are in use
* **Peer Address** is the real Internet address of this peer
* **Proxy ID** contains an ID of a tunnel created on a forwarder


