package main

import (
	"client/servant"
	"client/utils"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"time"

	"net/url"

	"github.com/gorilla/websocket"

	. "gfx/app/api/mdm"

	"gfx/library/service"

	"github.com/gogf/gf/crypto/gcrc32"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/guuid"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gproc"
	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
)

//BuildID 编译版本号
var BuildID = "0"

//写入文件一行,自动附加换行符
func WriteWithIoutil(name, content string) {
	//data := []byte(content)
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	if err == nil {
		f.WriteString(content)
	}
}

func ReadWithIoutil(name string) string {
	//data := []byte(content)
	f, err := os.Open(name)
	defer f.Close()
	if err == nil {
		fd, err := ioutil.ReadAll(f)
		if err != nil {
			g.Log().Errorf("read to fd fail %v", err)
			return ""
		}
		return string(fd)
	}
	return ""
}

func main() {

	sValue := os.Getenv("OVERSEER_BIN_CHECK")
	if sValue != "" {
		fmt.Printf(sValue)
		return
	}

	bOvrDbg := g.Cfg().GetBool("overseer.debug", false)
	sOvrURL := g.Cfg().GetString("overseer.url")
	nOvrIvl := g.Cfg().GetInt64("overseer.interval", 30)
	sMdmURL := g.Cfg().GetString("mdm.server")

	if sOvrURL == "" {
		g.Log().Fatal("overseer.url not defined!")
	}

	if sMdmURL == "" {
		g.Log().Fatal("mdm.server not defined!")
	}

	var bIsSlaveProcess = false
	if os.Getenv("OVERSEER_IS_SLAVE") == "1" {
		bIsSlaveProcess = true
	}
	g.Log().Infof("Main Version:%s IsChild:%v IsSlave:%v PPidOS:%v", BuildID, gproc.IsChild(), bIsSlaveProcess, gproc.PPidOS())

	overseer.Run(overseer.Config{
		Program:   prog,
		NoRestart: false,
		Fetcher: &fetcher.HTTP{
			URL:      sOvrURL,
			Interval: time.Duration(nOvrIvl) * time.Second,
		},
		//Fetcher: &fetcher.File{Path: "client2"},
		Debug: bOvrDbg,
	})
}

//prog(state) runs in a child process
func prog(state overseer.State) {

	sMdmURL := g.Cfg().GetString("mdm.server")
	sServerList := strings.Split(sMdmURL, "/") //	http://192.168.61.1:2389/mdm

	var bIsSlaveProcess = false
	if os.Getenv("OVERSEER_IS_SLAVE") == "1" {
		bIsSlaveProcess = true
	}

	srv := service.SystemService{
		Name:        "SysAgent",
		DisplayName: "SysAgent",
		Description: "Clientside for SysAdmin.",
		MainLoop: func() {

			//主进程直接退出
			if !bIsSlaveProcess {
				return
			}

			//连接断开自动重试，保持进程活动状态
			for {

				g.Log().Notice("----------------Wait for server online.----------------")

				sWsServer := sServerList[2]
				sWsContext := sServerList[3]

				u := url.URL{Scheme: "ws", Host: sWsServer, Path: "/" + sWsContext}
				g.Log().Infof("connecting to %s", u.String())

				//尝试建立连接
				c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
				if err != nil {
					g.Log().Errorf("dial:%v", err)

					//连接失败等30秒重试
					time.Sleep(time.Second * 30)
					continue
				}

				doWork(c, sMdmURL)
			}
		},
	}

	srv.Run()
}

func onClose(code int, text string) error {
	g.Log().Error("Closed\n")
	return nil
}

func genClientId() string {

	var sUUID string
	sUUID = ReadWithIoutil("identity")
	if sUUID == "" {
		uuid, _ := guuid.NewUUID()
		sUUID = fmt.Sprintf("%d", gcrc32.Encrypt(uuid.String()))
		WriteWithIoutil("identity", sUUID)
	}
	return sUUID
}

func doWork(c *websocket.Conn, sMdmURL string) {

	sClientId := genClientId()

	//连接建立成功，设置关闭方法和退出执行过程自动关闭连接
	c.SetCloseHandler(onClose)
	defer c.Close()

	var req = Request{}
	var token string
	for {
		errRdr := c.ReadJSON(&req)

		//读取信息失败，认为连接异常，自动退出当前连接
		if errRdr != nil {
			g.Log().Errorf("read:%v", errRdr)
			break
		}

		token = req.TraceId
		g.Log().Infof("recv %v", req)

		var sResult string
		var err error = nil
		var sReport string = ""

		if req.Cmd == "init" {
			sResult = fmt.Sprintf("%s|%s|%s", sClientId, BuildID, utils.GetSysInfo())
			sReport = fmt.Sprintf(`{"cmd":"%s","token":"%s","result":"%s","error":"%v"}`,
				req.Cmd, token, gbase64.EncodeString(sResult), "")
		} else {
			sResult, err = servant.ShellExec(req.Cmd)

			var sError string
			if err != nil {
				sError = err.Error()
			}

			sReport = fmt.Sprintf(`{"client":"%s","cmd":"%s","token":"%s","result":"%s","error":"%v"}`,
				sClientId, req.Cmd, token, gbase64.EncodeString(sResult), sError)
		}

		Report(sMdmURL+"/report", sReport)

	}
}
