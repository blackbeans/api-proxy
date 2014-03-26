package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

/**
 *用于国内访问苹果收据的代理
 *
 */
import _ "net/http/pprof"

const REQ_URL = "https://buy.itunes.apple.com/verifyReceipt"
const SD_REQ_URL = "https://sandbox.itunes.apple.com/verifyReceipt"

type RequestParams struct {
	Location   string
	RegionName string
}

type ResponseInfo struct {
	StatusCode int
	resultl    string
}

const (
	TIMEOUT = ""
)

func handleAppleADRecipt(rw http.ResponseWriter, r *http.Request) {
	innerHandle(rw, r, SD_REQ_URL)
}

func handleAppleRecipt(rw http.ResponseWriter, r *http.Request) {
	innerHandle(rw, r, REQ_URL)
}

func innerHandle(rw http.ResponseWriter, r *http.Request, url string) {
	data, err := ioutil.ReadAll(r.Body)
	if nil != err {
		log.Println("gen apple json args fail ", err.Error())
		return
	}

	defer r.Body.Close()
	total++
	respCh := make(chan []byte, 1)
	timeout := make(chan bool, 1)
	go func() {
		resp, err := http.Post(url, "application/json", bytes.NewReader(data))

		if nil != err {
			log.Printf("request apple %s ,%s : ", url, err.Error())

		} else {
			body, err := ioutil.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK {
				if nil != err {
					log.Println("recieve apple api body succ ,url : ", url, err)
					data, _ = json.Marshal(&ResponseInfo{500, "recieve response fail !"})
				} else {
					data = []byte(body)
					succ++
				}
			} else {
				data, _ = json.Marshal(&ResponseInfo{500, "r"})
			}
			defer resp.Body.Close()
		}
		if nil != data {
			respCh <- data
		}
		timeout <- false
	}()

	//模仿定时请求 10s瞪大
	select {
	case <-timeout:
	case <-time.After(10 * time.Second):
		data, _ = json.Marshal(&ResponseInfo{500, "recieve response timeout !"})
	}
	defer close(respCh)
	defer close(timeout)

	rw.Header().Set("contentType", "application/json")
	rw.Header().Add("charset", "utf-8")
	rw.Write(data)
}

var total int64 = 0
var succ int64 = 0

func main() {

	port := flag.String("port", ":7071", "-port=:7071 来定义端口!")
	flag.Parse()
	http.HandleFunc("/verifyReceipt", handleAppleRecipt)
	http.HandleFunc("/verifyReceipt/sandbox", handleAppleADRecipt)
	go func() {
		log.Fatal(http.ListenAndServe(*port, nil))
	}()

	var lastSucc int64 = 0
	var lastTotal int64 = 0
	//
	//模仿定时请求
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:

			go func() {
				log.Printf("%d/%d\t%d/%d", lastSucc, succ, lastTotal, total)
				lastSucc = succ
				lastTotal = total
			}()
		}
	}

}
