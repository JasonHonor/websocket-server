package mdm

import (
	"fmt"
	"os"

	"github.com/gogf/gf/encoding/gbase64"
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

		fmt.Printf("RequetWorker %v\n", wsWorker)

		// fetch job
		work := Job{serload: wsWorker.Index}
		JobQueue <- work

	} else {
		g.Log().Error(err)
		return
	}

	fmt.Printf("+")
}

//Report 接收客户端执行结果
func (c *HttpEntry) Report(r *ghttp.Request) {

	jsonData, err := r.GetJson()
	if err != nil {
		g.Log().Printf("解析Report请求失败:%s", err.Error())
		return
	}
	g.Log().Printf("Report %v\n", jsonData.MustToJsonString())

	connId := jsonData.Get("token")
	if mapByWorkerId.Get(connId) == nil {
		r.Response.WriteStatusExit(404, "")
		return
	}

	//继续处理客户端执行结果
	sBase64Ret := jsonData.Get("result").(string)

	if len(sBase64Ret) > 0 {
		sResult, err := gbase64.DecodeString(sBase64Ret)
		if err != nil {
			g.Log().Printf("ExecuteResult Error:%v\n", err.Error())
		} else {
			g.Log().Printf("ExecuteResult:%v\n", string(sResult))
		}
	}
}

func (c *HttpEntry) getOS(r *ghttp.Request) string {
	var sOs string = ""
	vOs := r.GetRouterVar("os")
	if vOs != nil {
		sOs = vOs.String()
	}

	if sOs == "" {
		sOs = "linux"
	}
	return sOs
}

//Upgrade 下载升级文件
func (c *HttpEntry) Upgrade(r *ghttp.Request) {

	sOs := c.getOS(r)

	sFile := g.Cfg().GetString("overseer." + sOs + ".source")

	if !FileExist(sFile) {
		g.Log().Debugf("Upgrading Failed %s", sFile)
		r.Response.WriteStatusExit(404)
		return
	}

	r.Response.ServeFileDownload(sFile, "client")
}

//Config 下载配置文件
func (c *HttpEntry) Config(r *ghttp.Request) {

	sOs := c.getOS(r)

	sFile := g.Cfg().GetString("overseer." + sOs + ".config")

	if !FileExist(sFile) {
		g.Log().Debugf("GetConfig Failed %s", sFile)
		r.Response.WriteStatusExit(404)
		return
	}

	r.Response.ServeFileDownload(sFile, "config.toml")
}

//Deploy 下载部署脚本，仅限Linux
func (c *HttpEntry) Deploy(r *ghttp.Request) {

	sServer := g.Cfg().GetString("overseer.server")

	g.Log().Debugf("Deploy by %s", sServer)

	r.Response.WriteTpl("deploy.html", g.Map{"server": sServer})
}

//FileExist 判断文件是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
