// Copyright 2013, Cong Ding. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Cong Ding <dinggnu@gmail.com>

package stun

import (
	"errors"
	"fmt"
	"net"
)

// Padding the length of the byte slice to multiple of 4.
func padding(bytes []byte) []byte {
	length := uint16(len(bytes))
	return append(bytes, make([]byte, align(length)-length)...)
}

// Align the uint16 number to the smallest multiple of 4, which is larger than
// or equal to the uint16 number.
func align(n uint16) uint16 {
	return (n + 3) & 0xfffc
}

func sendBindingReq(serverAddr string, c *Client) (*packet, string, error) {
	var conn net.Conn
	var err error
	var saddr *net.UDPAddr
	var laddr *net.UDPAddr
	var flag bool = true

	laddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", c.LocalAddr, c.LocalPort))
	flag = flag && (err == nil)
	saddr, err = net.ResolveUDPAddr("udp", serverAddr)
	flag = flag && (err == nil)

	if !flag {
		conn, err = net.Dial("udp", serverAddr)
		if err != nil {
			return nil, "", err
		}
	} else {
		conn, err = net.DialUDP("udp", laddr, saddr)
		if err != nil {
			return nil, "", err
		}
	}
	// Construct packet.
	packet := newPacket()
	packet.types = type_BINDING_REQUEST
	attribute := newSoftwareAttribute(packet, DefaultSoftwareName)
	packet.addAttribute(*attribute)
	attribute = newFingerprintAttribute(packet)
	packet.addAttribute(*attribute)
	// Send packet.
	localAddr := conn.LocalAddr().String()
	packet, err = packet.send(conn)
	if err != nil {
		return nil, "", err
	}
	err = conn.Close()
	return packet, localAddr, err
}

func sendChangeReq(serverAddr string, changeIP bool, changePort bool, c *Client) (*packet, error) {
	var conn net.Conn
	var err error
	var saddr *net.UDPAddr
	var laddr *net.UDPAddr
	var flag bool = true

	laddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", c.LocalAddr, c.LocalPort))
	flag = flag && (err == nil)
	saddr, err = net.ResolveUDPAddr("udp", serverAddr)
	flag = flag && (err == nil)

	if !flag {
		conn, err = net.Dial("udp", serverAddr)
		if err != nil {
			return nil, err
		}
	} else {
		conn, err = net.DialUDP("udp", laddr, saddr)
		if err != nil {
			return nil, err
		}
	}

	// Construct packet.
	packet := newPacket()
	packet.types = type_BINDING_REQUEST
	attribute := newSoftwareAttribute(packet, DefaultSoftwareName)
	packet.addAttribute(*attribute)
	attribute = newChangeReqAttribute(packet, changeIP, changePort)
	packet.addAttribute(*attribute)
	attribute = newFingerprintAttribute(packet)
	packet.addAttribute(*attribute)
	// Send packet.
	packet, err = packet.send(conn)
	if err != nil {
		return nil, err
	}
	err = conn.Close()
	return packet, err
}

func test1(serverAddr string, c *Client) (*packet, string, bool, *Host, error) {
	packet, localAddr, err := sendBindingReq(serverAddr, c)
	if err != nil {
		return nil, "", false, nil, err
	}
	if packet == nil {
		return nil, "", false, nil, nil
	}

	// RFC 3489 doesn't require the server return XOR mapped address.
	hostMappedAddr := packet.xorMappedAddr()
	if hostMappedAddr == nil {
		hostMappedAddr = packet.mappedAddr()
		if hostMappedAddr == nil {
			return nil, "", false, nil, errors.New("No mapped address.")
		}
	}

	hostChangedAddr := packet.changedAddr()
	if hostChangedAddr == nil {
		return nil, "", false, nil, errors.New("No changed address.")
	}
	changeAddr := hostChangedAddr.TransportAddr()
	identical := localAddr == hostMappedAddr.TransportAddr()
	return packet, changeAddr, identical, hostMappedAddr, nil
}

func test2(serverAddr string, c *Client) (*packet, error) {
	return sendChangeReq(serverAddr, true, true, c)
}

func test3(serverAddr string, c *Client) (*packet, error) {
	return sendChangeReq(serverAddr, false, true, c)
}

// Follow RFC 3489 and RFC 5389.
// Figure 2: Flow for type discovery process (from RFC 3489).
//                        +--------+
//                        |  Test  |
//                        |   I    |
//                        +--------+
//                             |
//                             |
//                             V
//                            /\              /\
//                         N /  \ Y          /  \ Y             +--------+
//          UDP     <-------/Resp\--------->/ IP \------------->|  Test  |
//          Blocked         \ ?  /          \Same/              |   II   |
//                           \  /            \? /               +--------+
//                            \/              \/                    |
//                                             | N                  |
//                                             |                    V
//                                             V                    /\
//                                         +--------+  Sym.      N /  \
//                                         |  Test  |  UDP    <---/Resp\
//                                         |   II   |  Firewall   \ ?  /
//                                         +--------+              \  /
//                                             |                    \/
//                                             V                     |Y
//                  /\                         /\                    |
//   Symmetric  N  /  \       +--------+   N  /  \                   V
//      NAT  <--- / IP \<-----|  Test  |<--- /Resp\               Open
//                \Same/      |   I    |     \ ?  /               Internet
//                 \? /       +--------+      \  /
//                  \/                         \/
//                  |Y                          |Y
//                  |                           |
//                  |                           V
//                  |                           Full
//                  |                           Cone
//                  V              /\
//              +--------+        /  \ Y
//              |  Test  |------>/Resp\---->Restricted
//              |   III  |       \ ?  /
//              +--------+        \  /
//                                 \/
//                                  |N
//                                  |       Port
//                                  +------>Restricted
func discover(c *Client) (NATType, *Host, error) {
	packet, changeAddr, identical, host, err := test1(c.ServerAddr(), c)
	if err != nil {
		return NAT_ERROR, nil, err
	}
	if packet == nil {
		return NAT_BLOCKED, nil, nil
	}
	packet, err = test2(c.ServerAddr(), c)
	if err != nil {
		return NAT_ERROR, host, err
	}
	if identical {
		if packet == nil {
			return NAT_SYMETRIC_UDP_FIREWALL, host, nil
		}
		return NAT_NONE, host, nil
	}
	if packet != nil {
		return NAT_FULL, host, nil
	}
	packet, _, identical, _, err = test1(changeAddr, c)
	if err != nil {
		return NAT_ERROR, host, err
	}
	if packet == nil {
		// It should be NAT_BLOCKED, but will be detected in the first
		// step. So this will never happen.
		return NAT_UNKNOWN, host, nil
	}
	if identical {
		packet, err = test3(c.ServerAddr(), c)
		if err != nil {
			return NAT_ERROR, host, err
		}
		if packet == nil {
			return NAT_PORT_RESTRICTED, host, nil
		}
		return NAT_RESTRICTED, host, nil
	}
	return NAT_SYMETRIC, host, nil
}
