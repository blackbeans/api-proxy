package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const REQ_URL = "https://maps.googleapis.com/maps/api/place/nearbysearch/json?" +
	"location={0}&radius=2000&sensor=false&name={1}&key={2}"

type RequestParams struct {
	Location   string
	RegionName string
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
			data = []byte("{\"result\",\"request google fail \"}")
		} else {
			body, err := ioutil.ReadAll(result.Body)
			if nil != err {
				log.Println("recieve google api body fail ,url : ", reqURL, err)
				data = []byte("{\"result\",false}")
			} else {
				data = []byte(body)
			}
		}
	}
	rw.Header().Set("contentType", "text/json")
	rw.Header().Add("charset", "utf-8")
	rw.Write(data)
}

func main() {
	http.HandleFunc("/api/nearby", handleGoogleApi)

	http.ListenAndServe(":7070", nil)
}
