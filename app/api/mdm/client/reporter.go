package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gogf/gf/frame/g"
)

func Report(url string, data string) {
	g.Log().Printf("DoRport.%s %s\n", url, data)

	btData := []byte(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(btData))
	if err != nil {
		g.Log().Error("Error reading request. ", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	//eq.Header.Set("Host", "httpbin.org")

	// Create and Add cookie to request
	//cookie := http.Cookie{Name: "cookie_name", Value: "cookie_value"}
	//req.AddCookie(&cookie)

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}
	defer client.CloseIdleConnections()

	// Validate cookie and headers are attached
	g.Log().Println(req.Cookies())
	g.Log().Println(req.Header)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		g.Log().Error("Error reading response. ", err)
		return
	}
	defer resp.Body.Close()

	g.Log().Println("response Status:", resp.Status)
	g.Log().Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		g.Log().Error("Error reading body. ", err)
		return
	}

	g.Log().Printf("%s\n", body)
}
