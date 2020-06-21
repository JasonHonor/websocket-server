package mdm

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type HttpEntry struct {
}

func (c *HttpEntry) Index(r *ghttp.Request) {

	// 初始化WebSocket请求
	if ws, err := r.WebSocket(); err == nil {

		wsWorker := &WsWorker{
			Websocket: ws,
			Request:   r,
		}

		// fetch job
		work := Job{serload: wsWorker.Index}
		JobQueue <- work

	} else {
		g.Log().Error(err)
		return
	}

	fmt.Printf("+")
}

func (c *HttpEntry) Report(r *ghttp.Request) {
	jsonData, err := r.GetJson()
	if err != nil {
		g.Log().Printf("解析Report请求失败:%s", err.Error())
		return
	}
	g.Log().Printf("Report %v\n", jsonData.MustToJsonString())
}
