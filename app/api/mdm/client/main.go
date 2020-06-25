package main

import (
	"client/servant"
	"fmt"
	"os"
	"path/filepath"

	"time"

	"net/url"

	"github.com/gorilla/websocket"

	. "gfx/app/api/mdm"

	"gfx/library/service"

	"github.com/gogf/gf/encoding/gbase64"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gproc"
	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
)

var BuildID = "0"

func getAppDir() string {
	dir, errDir := filepath.Abs(filepath.Dir(os.Args[0]))
	if errDir != nil {
		return ""
	}

	return dir
}

func main() {

	g.Cfg().SetPath(getAppDir())

	bOvrDbg := g.Cfg().GetBool("overseer.debug", false)
	sOvrUrl := g.Cfg().GetString("overseer.url")
	nOvrIvl := g.Cfg().GetInt64("overseer.interval", 30)

	if sOvrUrl == "" {
		g.Log().Fatal("overseer.url not defined!")
	}

	overseer.Run(overseer.Config{
		Program:   prog,
		NoRestart: false,
		Fetcher: &fetcher.HTTP{
			URL:      sOvrUrl,
			Interval: time.Duration(nOvrIvl) * time.Second,
		},
		//Fetcher: &fetcher.File{Path: "client2"},
		Debug: bOvrDbg,
	})
	//prog(overseer.DisabledState)
}

//prog(state) runs in a child process
func prog(state overseer.State) {

	var bIsSlaveProcess = false
	if os.Getenv("OVERSEER_IS_SLAVE") == "1" {
		bIsSlaveProcess = true
	}

	g.Log().Infof("Version:%s IsChild:%v IsSlave:%v\n", BuildID, gproc.IsChild(), bIsSlaveProcess)

	srv := service.SystemService{
		Name:        "SysAgent",
		DisplayName: "SysAgent",
		Description: "Clientside for SysAdmin.",
		MainLoop: func() {

			if !bIsSlaveProcess {
				return
			}

			g.Log().Infof("I'm Running! IsChild:%v\t", gproc.IsChild())

			u := url.URL{Scheme: "ws", Host: "127.0.0.1:8199", Path: "/mdm"}
			g.Log().Infof("connecting to %s", u.String())

			c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				g.Log().Errorf("dial:", err)
				return
			}

			c.SetCloseHandler(onClose)
			defer c.Close()

			var req = Request{}
			var token string
			for {
				errRdr := c.ReadJSON(&req)
				if errRdr != nil {
					g.Log().Errorf("read:%v", errRdr)
					return
				}

				token = req.TraceId
				g.Log().Infof("recv %v\n", req)

				s, err := servant.ShellExec(req.Cmd)

				var sError string
				if err != nil {
					sError = err.Error()
				}

				Report("http://127.0.0.1:8199/mdm/report", fmt.Sprintf(`{"cmd":"%s","token":"%s","result":"%s","error":"%v"}`,
					req.Cmd, token,
					gbase64.EncodeString(s), gbase64.EncodeString(sError)))

			}
		},
	}

	srv.Run()
}

func onClose(code int, text string) error {
	fmt.Printf("Closed\n")
	return nil
}
