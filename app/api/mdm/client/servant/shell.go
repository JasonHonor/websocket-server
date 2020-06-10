package servant

import (
	"github.com/gogf/gf/os/gproc"
)

func ShellExec(cmd string) (string, error) {
	r, err := gproc.ShellExec(cmd)
	return r, err
}
