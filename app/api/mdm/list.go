package mdm

import (
	"fmt"

	"github.com/gogf/gf/net/ghttp"
)

//列出所有的客户端信息
func (c *HttpEntry) List(r *ghttp.Request) {
	r.Response.Write(fmt.Sprintf("clients count %v queue:%v", mapByWorker.Size(), len(JobQueue)))
}

func (c *HttpEntry) Push(r *ghttp.Request) {
	for _, workerId := range mapByWorkerId.Keys() {
		worker := mapByWorkerId.Get(workerId)
		if worker != nil {
			wsWorker := worker.(*WsWorker)
			PushJson(wsWorker.Websocket, "date", wsWorker.WorkerId)
		}
	}
}
