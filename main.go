package main

import (
	_ "gfx/boot"

	"gfx/library/service"

	"github.com/gogf/gf/frame/g"
)

func main() {

	svc := service.SystemService{
		Name:        "SysCenter",
		DisplayName: "SysCenter",
		Description: "Centerside for syscenter.",
		MainLoop: func() {
			g.Log().SetStack(false)
			g.Server().Run()
		},
	}

	svc.Run()
}
