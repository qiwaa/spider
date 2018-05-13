package main

import (
	"github.com/chenminjian/spider/krpc"
	"github.com/chenminjian/spider/utils"
)

func main() {

	utils.SetProcessName("spider")

	mgr := krpc.NewKrpcManager()
	go mgr.Run()
	mgr.Join()

	select {}
}
