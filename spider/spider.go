package main

import (
	"github.com/chenminjian/spider/krpc"
	"github.com/chenminjian/spider/utils"
	"github.com/chenminjian/spider/controller"
)

func main() {

	utils.SetProcessName("spider")

	mgr := krpc.Init()
	go mgr.Run()
	mgr.Join()

	controller.Run()

	select {}
}
