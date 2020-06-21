package main

import (
	"client/servant"
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/jvehent/service-go"

	. "gfx/app/api/mdm"

	"github.com/gogf/gf/encoding/gbase64"
)

var log service.Logger

func main() {
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
	log.Info("I'm Running!")

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8199", Path: "/mdm"}
	log.Info("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error("dial:", err)
		return
	}

	c.SetCloseHandler(onClose)
	defer c.Close()

	var resp = Response{}
	var req = Request{}
	var traceId string
	for {
		errRdr := c.ReadJSON(&req)
		if errRdr != nil {
			log.Info("read:", errRdr)
			return
		}

		traceId = req.TraceId
		fmt.Printf("recv %v\n", req)

		s, err := servant.ShellExec(req.Cmd)

		c.WriteMessage(websocket.TextMessage,
			[]byte(fmt.Sprintf(`{"cmd":"%s","trace_id":"%s","result":"%s","error":"%v"}`, req.Cmd, traceId,
				gbase64.EncodeString(s), err)))

		errRdr2 := c.ReadJSON(&resp)
		if errRdr2 != nil {
			log.Error("read:", errRdr2)
			return
		}
		fmt.Printf("recv2 %v\n", resp)
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
