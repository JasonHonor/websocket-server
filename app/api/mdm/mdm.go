package mdm

import (
	"github.com/gogf/gf/net/ghttp"
)

type WsWorker struct {
	Request   *ghttp.Request
	Websocket *ghttp.WebSocket
}

type Request struct {
	Cmd     string      `json:"cmd" v:"cmd@required#消息类型不能为空"`
	TraceId string      `json:"trace_id" v:"trace_id@required#消息编号不能为空"`
	Param   interface{} `json:"cmd_param"`
}

type Response struct {
	Cmd     string      `json:"cmd"`
	TraceId string      `json:"trace_id"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
