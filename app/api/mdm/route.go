package mdm

import (
	"fmt"

	"github.com/gogf/gf/container/gmap"
)

type WSHandler func(context *WsWorker, request *Request, response *Response)

var (
	gRoutes = gmap.NewStrAnyMap(true)
)

//获取命令路由
func GetRoute(r Request) (WSHandler, error) {

	if !gRoutes.Contains(r.Cmd) {
		return nil, fmt.Errorf("Handler for %s not found!", r.Cmd)
	}

	return gRoutes.Get(r.Cmd).(WSHandler), nil
}

//设置命令路由
func SetRoute(cmd string, handler WSHandler) {
	gRoutes.Set(cmd, handler)
}
