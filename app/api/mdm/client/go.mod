module client

go 1.14

require (
	gfx v0.0.0-00010101000000-000000000000
	github.com/gogf/gf v1.13.1
	github.com/gorilla/websocket v1.4.2
	github.com/jpillora/overseer v1.1.4
	github.com/judwhite/go-svc v1.1.2
	github.com/jvehent/service-go v0.0.0-20160824215813-0da6d786ded5
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	golang.org/x/sys v0.0.0-20200610111108-226ff32320da // indirect
)

replace gfx => ../../../../../websocket-server

replace github.com/jpillora/overseer => /home/wangjinhui/projects/overseer
