package api

import (
	"net/http"
	"github.com/chenminjian/spider/krpc"
	"encoding/json"
)

func Metric(resp http.ResponseWriter, req *http.Request) {

	metric := krpc.Mgr.Metric
	buf, _ := json.Marshal(metric)

	resp.Write(buf)
}
