package mdm

import (
	"fmt"

	"github.com/gogf/gf/net/ghttp"
)

//列出所有的客户端信息
func (c *HttpEntry) List(r *ghttp.Request) {

	r.Response.Write(fmt.Sprintf("clients count %v queue:%v", devices.Size(), len(JobQueue)))
}
