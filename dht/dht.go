package dht

import (
	"github.com/chenminjian/spider/route"
	"github.com/chenminjian/spider/config"
	"net"
	"github.com/chenminjian/spider/handle"
	"github.com/chenminjian/spider/utils"
	"github.com/chenminjian/spider/container"
	"fmt"
)

type Dht struct {
	*config.Config                    // dht config
	Node         *config.Node         // this dht's node
	Conn         *net.UDPConn         // udp server
	RoutingTable *route.Routingtable  // route table
	packets      chan *config.Packet  // receive packet chan
	handler      handle.MsgHandler    // response handler
	BlackList    *container.BlackList // black node list
}

func NewDht() *Dht {
	conf := config.NewFixedConfig()
	node, _ := config.NewNode(utils.RandomString(20), conf.Addr)
	return &Dht{
		Config:       conf,
		Node:         node,
		RoutingTable: route.NewRoutingTable(),
		packets:      make(chan *config.Packet, 1024),
		BlackList:    container.NewBlackList(1024),
	}
}

func (d *Dht) SetHandler(handler handle.MsgHandler) {
	d.handler = handler
}

func (d *Dht) init() {
	listener, err := net.ListenPacket("udp", d.Addr)
	if err != nil {
		panic(err)
	}
	d.Conn = listener.(*net.UDPConn)

}

// receives message from udp.
func (d *Dht) listen() {
	go func() {
		buff := make([]byte, 8192)
		for {
			n, raddr, err := d.Conn.ReadFromUDP(buff)
			if err != nil {
				fmt.Println("ReadFromUDP error!")
				continue
			}
			fmt.Println("receive packet!")

			d.packets <- &config.Packet{
				Data: buff[:n],
				Addr: raddr,
			}

			fmt.Println("after send packet to chan!")
		}
	}()
}

func (d *Dht) Run() {
	d.init()
	d.listen()

	var p *config.Packet

	for {
		select {
		case p = <-d.packets:
			fmt.Println("here!")
			go d.handler.Handle(p)
		}
	}
}
