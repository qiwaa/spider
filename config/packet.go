package config

import "net"

type Packet struct {
	Data []byte
	Addr *net.UDPAddr
}
