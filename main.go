package main

import (
	_ "gfx/boot"
	"os"
	"path/filepath"

	"gfx/library/service"

	"github.com/gogf/gf/frame/g"
)

func getAppDir() string {
	dir, errDir := filepath.Abs(filepath.Dir(os.Args[0]))
	if errDir != nil {
		return ""
	}

	return dir
}

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
