package boot

import (
	"os"
	"path/filepath"

	"github.com/gogf/gf-swagger/swagger"
	"github.com/gogf/gf/frame/g"
)

func getAppDir() string {
	dir, errDir := filepath.Abs(filepath.Dir(os.Args[0]))
	if errDir != nil {
		return ""
	}

	return dir
}

// 用于应用初始化。
func init() {
	s := g.Server()
	s.Plugin(&swagger.Swagger{})
}
