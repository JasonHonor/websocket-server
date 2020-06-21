package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Report(url string, data string) {
	fmt.Printf("DoRport.%s %s\n", url, data)

	btData := []byte(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(btData))
	if err != nil {
		log.Error("Error reading request. ", err)
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

	// Validate cookie and headers are attached
	fmt.Println(req.Cookies())
	fmt.Println(req.Header)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error reading response. ", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading body. ", err)
		return
	}

	fmt.Printf("%s\n", body)
}
