package mdm

import (
	"fmt"

	"github.com/gogf/gf/net/ghttp"
)

//列出所有的客户端信息
func (c *HttpEntry) List(r *ghttp.Request) {
	r.Response.Write(fmt.Sprintf("clients count %v queue:%v", deviceMapBySocket.Size(), len(JobQueue)))
}

func (c *HttpEntry) Push(r *ghttp.Request) {
	for _, clientName := range deviceMapByName.Keys() {
		ws := deviceMapByName.Get(clientName)
		PushJson(ws.(*ghttp.WebSocket), "pwd", "1")
	}
}
