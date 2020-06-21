package mdm

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/crypto/gcrc32"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/util/guuid"
	"github.com/gogf/gf/util/gvalid"
)

var (
	// 使用默认的并发安全Map
	deviceMapBySocket = gmap.New(true)
	deviceMapByName   = gmap.New(true)

	// 使用并发安全的Set，用于设备唯一性校验
	names = gset.NewStrSet(true)

	// 使用特定的缓存对象，不使用全局缓存对象
	cache = gcache.New()
)

// @summary 设备管理通道
// @description 实现设备管理相关通讯。
// @tags
// @produce
// @router  /mdm [POST]
// @success 200 {string} string ""
func (c *WsWorker) Index() {

	var traceId string
	//r := c.Request
	ws := c.Websocket

	uuid, _ := guuid.NewUUID()
	sUUID := fmt.Sprintf("%d", gcrc32.Encrypt(uuid.String()))
	c.WriteJson(ws, 0, "", "init", sUUID)
	traceId = sUUID

	c.appendDevice(ws)

	for {
		// 阻塞读取WS数据
		_, msgByte, err := ws.ReadMessage()
		if err != nil {
			c.onError(err, "read", traceId)
			break
		}

		req := Request{}
		resp := Response{}

		// JSON参数解析
		if err := gjson.DecodeTo(msgByte, &req); err != nil {
			c.onError(err, "decode", traceId, "消息格式不正确")
			continue
		}
		// 数据校验
		if e := gvalid.CheckStruct(req, nil); e != nil {
			c.onError(e, req.Cmd, traceId, "消息验证失败")
			continue
		}

		traceId = req.TraceId

		// 如果失败，那么表示断开，这里清除用户信息
		handler, errHandler := GetRoute(req)

		if errHandler != nil {
			c.onError(errHandler, req.Cmd, traceId)
		} else {
			//执行自定义处理过程
			handler(c, &req, &resp)
		}
	}
	ws.Close()
}

func (c *WsWorker) appendDevice(r *ghttp.WebSocket) {
	//保存在线信息
	name := r.RemoteAddr().String()
	names.Add(name)

	deviceMapByName.Set(name, c.Websocket)
	deviceMapBySocket.Set(c.Websocket, name)

	r.SetCloseHandler(c.onClose)
}

func (c *WsWorker) removeDevice(r *ghttp.WebSocket) {

	//保存在线信息
	name := r.RemoteAddr().String()
	names.Remove(name)
	deviceMapBySocket.Remove(r)
	deviceMapByName.Remove(name)
}

func (c *WsWorker) onClose(int, string) error {
	fmt.Printf("OnClosed\n")
	c.removeDevice(c.Websocket)
	return nil
}

func (c *WsWorker) buildError(errHandler error, ctx ...string) string {

	g.Log().Error(errHandler)
	// 使用特定的缓存对象，不使用全局缓存对象

	var msg string
	if len(ctx) > 0 {
		msg = fmt.Sprintf("[%s]%v", ctx[0], errHandler.Error())
	} else {
		msg = fmt.Sprintf("%v", errHandler.Error())
	}
	return msg
}

func (c *WsWorker) onError(errHandler error, cmd, traceId string, ctx ...string) {

	if errHandler == nil {
		return
	}

	sErr := errHandler.Error()

	fmt.Printf("OnError %v\n", errHandler)

	//处理连接意外断开的情况
	if strings.Contains(sErr, "websocket: close 1006") ||
		strings.Contains(sErr, "use of closed network connection") ||
		strings.Contains(sErr, "connection reset by peer") ||
		strings.Contains(sErr, "broken pipe") {
		c.onClose(1, sErr)
		return
	}

	c.WriteJson(c.Websocket, 1, c.buildError(errHandler, ctx...), cmd, traceId)
}

func (c *WsWorker) onErrorExit(errHandler error, cmd, traceId string, ctx ...string) {
	c.WriteJsonExit(c.Request, c.Websocket, 1, c.buildError(errHandler, ctx...), cmd, traceId)
}

// 标准返回结果数据结构封装。
func (c *WsWorker) WriteJson(ws *ghttp.WebSocket, code int, message, cmd, traceId string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	err := ws.WriteJSON(Response{
		Code:    code,
		Message: message,
		Data:    responseData,
		Cmd:     cmd,
		TraceId: traceId,
	})

	if err != nil {
		g.Log().Error(err)
	}
}

// 返回JSON数据并退出当前HTTP执行函数。
func (c *WsWorker) WriteJsonExit(r *ghttp.Request, ws *ghttp.WebSocket, err int, msg, cmd, traceId string, data ...interface{}) {
	c.WriteJson(ws, err, msg, cmd, traceId, data...)
	r.Exit()
}

func PushJson(ws *ghttp.WebSocket, cmd, traceId string, data ...interface{}) {

	err := ws.WriteJSON(Request{
		Cmd:     cmd,
		TraceId: traceId,
		Param:   data,
	})

	if err != nil {
		g.Log().Error(err)
	}
}
