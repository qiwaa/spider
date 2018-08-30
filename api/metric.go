package api

import (
	"encoding/json"

	"github.com/chenminjian/spider/krpc"
	"github.com/gin-gonic/gin"
)

func Metric(c *gin.Context) {

	metric := krpc.Mgr.Metric
	buf, _ := json.Marshal(metric)

	c.JSON(200, gin.H{
		"code":    0,
		"message": "",
		"data":    string(buf),
	})
}
