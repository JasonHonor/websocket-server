package service

import (
	_ "gfx/boot"

	_ "gfx/router"

	"fmt"
	"os"

	//	"time"

	"github.com/jvehent/service-go"
)

var log service.Logger

type SystemService struct {
	Name        string
	DisplayName string
	Description string
	MainLoop    Mainfn
}

type Mainfn func()

var exit = make(chan struct{})

func (m *SystemService) Run() {
	var s, err = service.NewService(m.Name, m.DisplayName, m.Description)
	log = s

	if err != nil {
		fmt.Printf("%s unable to start: %s", m.DisplayName, err)
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
			fmt.Printf("Service \"%s\" installed.\n", m.DisplayName)
		case "remove":
			err = s.Remove()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", m.DisplayName)
		case "run":
			m.doWork()
		case "start":
			err = s.Start()
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" started.\n", m.DisplayName)
		case "stop":
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" stopped.\n", m.DisplayName)
		}
		return
	}
	err = s.Run(func() error {
		// start
		go m.doWork()
		return nil
	}, func() error {
		// stop
		m.stopWork()
		return nil
	})
	if err != nil {
		s.Error(err.Error())
	}
}

func (m *SystemService) doWork() {
	log.Info("I'm Running!")
	m.MainLoop()
}

func (m *SystemService) stopWork() {
	log.Info("I'm Stopping!")
	exit <- struct{}{}
}
