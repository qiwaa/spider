package handle

import "github.com/chenminjian/spider/config"

type MsgHandler interface {
	Handle(packet *config.Packet)
}
