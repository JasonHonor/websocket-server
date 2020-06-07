package main

import (
	"fmt"
	"os"

	"github.com/jvehent/service-go"
)

/*
//Remote Control Mode:
	//		串行化：receive(run-book)->execute(run-book)->report{send->receive}->nop(wait receive).
	//      并行：receive(thread1)				接收所有消息，区分job和response，job->newJob,response->workerThread
	//			 newJob(thread1)				只处理job任务
	//				PickJob(thread2)
	//				Execute(thread2)
	//				Report(thread2)
	//Execute Command Mode:
	//		send->receive->nop.

	//CommandList:
	//1. download
	//2. upgrade
	//3. shell
	//4. queue:队列  job参数：a.默认队列 default,b.后台线程  thread,c.新建队列 queue。
	//5. job参数： 任务编号：jobId, 客户端列表：agentList，命令:command,队列名称:queue(为空，默认队列)。
*/

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
	/*
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
				log.Info("Still running...")
			case <-exit:
				ticker.Stop()
				return
			}
		}
	*/

}

func stopWork() {
	log.Info("I'm Stopping!")
	exit <- struct{}{}
}
