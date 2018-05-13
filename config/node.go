package config

import (
	"net"
	"time"
	"errors"
	"github.com/chenminjian/spider/utils"
)

type Node struct {
	Id             string // 160 bit
	Addr           *net.UDPAddr
	LastActiveTime time.Time
}

func NewNode(id string, addr string) (*Node, error) {
	if len(id) != 20 {
		return nil, errors.New("node id illegal!")
	}

	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	return &Node{Id: id, Addr: udpAddr, LastActiveTime: time.Now()}, nil
}

func NewNodeFromCompactNodeInfo(compactNodeInfo string) (*Node, error) {
	if len(compactNodeInfo) != 26 {
		return nil, errors.New("compactNodeInfo' length should 26!")
	}

	id := compactNodeInfo[:20]
	ip, port, _ := decodeCompactIPPortInfo(compactNodeInfo[20:])

	return NewNode(id, utils.GenAddress(ip.String(), port))
}

func decodeCompactIPPortInfo(info string) (ip net.IP, port int, err error) {
	if len(info) != 6 {
		err = errors.New("compact info should be 6-length long")
		return
	}

	ip = net.IPv4(info[0], info[1], info[2], info[3])
	port = int((uint16(info[4]) << 8) | uint16(info[5]))
	return
}
