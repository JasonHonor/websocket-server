package main

import (
	"client/servant"
	"fmt"
	"net/url"
	"os"

	"time"

	"github.com/gorilla/websocket"
	"github.com/jvehent/service-go"

	. "gfx/app/api/mdm"

	"github.com/gogf/gf/encoding/gbase64"

	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"

	"github.com/gogf/gf/frame/g"
)

var log service.Logger

var BuildID = "0"

func main() {
	overseer.Run(overseer.Config{
		Program:   prog,
		NoRestart: false,
		Fetcher: &fetcher.HTTP{
			URL:      "http://127.0.0.1:8199/mdm/upgrade",
			Interval: 30 * time.Second,
		},
		//Fetcher: &fetcher.File{Path: "client2"},
		Debug: true,
	})
}

//prog(state) runs in a child process
func prog(state overseer.State) {
	fmt.Println("Version:\t", BuildID)

	var name = "SysAgent"
	var displayName = "SysAgent"
	var desc = "Agent for syscenter."

	var s, err = service.NewService(name, displayName, desc)
	log = s

	if err != nil {
		fmt.Printf("%s unable to start: %s", displayName, err)
		return
	}

	if len(os.Args) > 1 {
		var err error
		verb := os.Args[1]
		switch verb {
		case "install":
			err = s.Install()
			if err != nil {
				fmt.Printf("Failed to install: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" installed.\n", displayName)
		case "remove":
			err = s.Remove()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", displayName)
		case "run":
			doWork()
		case "start":
			err = s.Start()
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" started.\n", displayName)
		case "stop":
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" stopped.\n", displayName)
		}
		return
	}
	err = s.Run(func() error {
		// start
		go doWork()
		return nil
	}, func() error {
		// stop
		stopWork()
		return nil
	})
	if err != nil {
		s.Error(err.Error())
	}

}

var exit = make(chan struct{})

func doWork() {
	g.Log().Printf("I'm Running!")

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8199", Path: "/mdm"}
	g.Log().Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		g.Log().Printf("dial:", err)
		return
	}

	c.SetCloseHandler(onClose)
	defer c.Close()

	var req = Request{}
	var token string
	for {
		errRdr := c.ReadJSON(&req)
		if errRdr != nil {
			g.Log().Printf("read:%v", errRdr)
			return
		}

		token = req.TraceId
		g.Log().Printf("recv %v\n", req)

		s, err := servant.ShellExec(req.Cmd)

		var sError string
		if err != nil {
			sError = err.Error()
		}

		Report("http://127.0.0.1:8199/mdm/report", fmt.Sprintf(`{"cmd":"%s","token":"%s","result":"%s","error":"%v"}`,
			req.Cmd, token,
			gbase64.EncodeString(s), gbase64.EncodeString(sError)))
	}
}

func stopWork() {
	log.Info("I'm Stopping!")
	exit <- struct{}{}
}

func onClose(code int, text string) error {
	fmt.Printf("Closed\n")
	return nil
}
