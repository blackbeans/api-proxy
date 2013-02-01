package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

import _ "net/http/pprof"

const REQ_URL = "https://maps.googleapis.com/maps/api/place/nearbysearch/json?" +
	"location={0}&radius=2000&sensor=false&name={1}&key={2}"

type RequestParams struct {
	Location   string
	RegionName string
}

type ResponseInfo struct {
	StatusCode int
	resultl    string
}

func handleGoogleApi(rw http.ResponseWriter, r *http.Request) {
	location := r.FormValue("location")
	name := r.FormValue("name")
	licence := r.FormValue("token")
	var data []byte
	//参数验证
	if len(location) <= 0 && len(name) <= 0 && len(licence) <= 0 {
		data = []byte("{\"reuslt\":\"invalid params\"}")

	} else {

		var reqURL string = strings.Replace(REQ_URL, "{0}", location, -1)
		reqURL = strings.Replace(reqURL, "{1}", name, -1)
		reqURL = strings.Replace(reqURL, "{2}", licence, -1)
		log.Println("request url :", reqURL)
		result, err := http.Get(reqURL)
		defer result.Body.Close()

		if nil != err {
			log.Println("request google err ,url : ", reqURL, err)
			data, _ = json.Marshal(&ResponseInfo{500, "request google fail!"})
		} else {

			body, err := ioutil.ReadAll(result.Body)

			if result.StatusCode == http.StatusOK {
				if nil != err {
					log.Println("recieve google api body fail ,url : ", reqURL, err)
					data, _ = json.Marshal(&ResponseInfo{500, "recieve response fail !"})
				} else {
					data = []byte(body)
				}
			} else {
				data, _ = json.Marshal(&ResponseInfo{500, "r"})
			}

		}
	}
	rw.Header().Set("contentType", "text/json")
	rw.Header().Add("charset", "utf-8")
	rw.Write(data)
}

const (
	url = "http://ec2-54-248-164-29.ap-northeast-1.compute.amazonaws.com:7070/api/nearby?location=41.536497,123.601047&name=鑫雅轩小吃部&token=$youtoken"
)

func main() {

	http.HandleFunc("/api/nearby", handleGoogleApi)

	go func() {
		log.Fatal(http.ListenAndServe(":7070", nil))
	}()
	ch := make(chan time.Time)
	select {
	case ch <- time.Unix(0, time.Now()):
		log.Println(<-ch)

	}

	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			log.Println(time.Now())
			go func() {
				reps, err := http.Get(url)
				defer reps.Body.Close()
				if nil != err {
					return
				}

				data, _ := ioutil.ReadAll(reps.Body)
				log.Println(string(data))
			}()
		}
	}

}
