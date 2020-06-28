package router

import (
	"gfx/app/api/mdm"
	//"gfx/app/api/chat"
	//"gfx/app/api/curd"
	//"gfx/app/api/user"
	//	"gfx/app/service/middleware"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"

	_ "gfx/app/api/mdm"
)

// 你可以将路由注册放到一个文件中管理，
// 也可以按照模块拆分到不同的文件中管理，
// 但统一都放到router目录下。
func init() {
	s := g.Server()

	// 某些浏览器直接请求favicon.ico文件，特别是产生404时
	s.SetRewrite("/favicon.ico", "/resource/image/favicon.ico")

	// 分组路由注册方式
	s.Group("/", func(group *ghttp.RouterGroup) {

		//group.Middleware(middleware.CORS)
		entry := new(mdm.HttpEntry)
		group.ALL("/mdm/config/{os}", entry.Config)
		group.ALL("/mdm/upgrade/{os}", entry.Upgrade)
		//group.ALL("/mdm/deplpy", entry.Deploy)
		//group.ALL("/mdm/list", entry.List)
		//group.ALL("/mdm/push", entry.Push)
		group.ALL("/mdm", entry)
	})
}
