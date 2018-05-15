package controller

import (
	"net/http"
	"github.com/chenminjian/spider/api"
)

func Run() {
	http.HandleFunc("/metric", api.Metric)
	http.ListenAndServe(":8000", nil)
}
