package main

import (
	"net/http"
	"time"

	"github.com/chenminjian/spider/api"
	"github.com/chenminjian/spider/conf"
	"github.com/chenminjian/spider/krpc"
	"github.com/chenminjian/spider/model/dao"
	"github.com/chenminjian/spider/utils"
)

func main() {

	utils.SetProcessName("spider")

	if err := conf.Init("conf/config.toml"); err != nil {
		panic(err)
	}

	dao.Connect()

	mgr := krpc.Init()
	go mgr.Run()
	mgr.Join()

	r := api.Init()

	s := &http.Server{
		Addr:           ":8000",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   31 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
