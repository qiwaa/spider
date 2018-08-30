package handle

import "github.com/chenminjian/spider/config"

// MsgHandler is a handler which handles dht node's response and query.
type MsgHandler interface {
	Handle(packet *config.Packet)
}
