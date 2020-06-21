package mdm

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var wg sync.WaitGroup

var nConns = 1000
var nActions = 1
var nWait = 0
var bWriteFirstMode = false

func main() {
	host := os.Args[1]
	code := os.Args[2]

	testConnect(host, code, 1, bWriteFirstMode, nWait)
}

func TestConnect(t *testing.T) {

	if bWriteFirstMode == false {
		nActions = 1
	}

	for i := 0; i < nConns; i++ {

		wg.Add(1)
		go testConnect("127.0.0.1:8199", "", 0, bWriteFirstMode, nWait)
	}

	wg.Wait()
}

func onClose(code int, text string) error {
	fmt.Printf("Closed\n")
	wg.Done()

	return nil
}

func testConnect(host, code string, cntx int, writeFirstMode bool, nWait int) {

	defer wg.Done()

	u := url.URL{Scheme: "ws", Host: host, Path: "/mdm"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	c.SetCloseHandler(onClose)
	defer c.Close()

	var resp = Response{}

	var traceId string

	for i := 0; i < nActions; i++ {

		if writeFirstMode { //客户端请求模式,REQ-RESP
			c.WriteMessage(websocket.TextMessage,
				[]byte(fmt.Sprintf(`{"cmd":"get_session_id","trace_id":"%s","term_no":"rk3288|TESTWS"}`, traceId)))

			errRdr := c.ReadJSON(&resp)
			if errRdr != nil {
				log.Println("read:", errRdr)
				return
			}

			fmt.Printf("recv %v\n", resp)
			traceId = resp.TraceId
		} else { //推送处理模式,连接后获得会话ID. CONN->REQ->RESP
			errRdr := c.ReadJSON(&resp)
			if errRdr != nil {
				log.Println("read:", errRdr)
				return
			}

			traceId = resp.TraceId
			fmt.Printf("recv %v\n", resp)

			c.WriteMessage(websocket.TextMessage,
				[]byte(fmt.Sprintf(`{"cmd":"get_server_time","trace_id":"%s","term_no":"rk3288|TESTWS"}`, traceId)))

			errRdr2 := c.ReadJSON(&resp)
			if errRdr2 != nil {
				log.Println("read:", errRdr2)
				return
			}
			fmt.Printf("recv2 %v\n", resp)
		}

		time.Sleep(time.Duration(nWait) * time.Millisecond)
	}
}
