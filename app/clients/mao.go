package clients

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type MAOParam struct {
	MAOHost string
	APIKey  string
}

func (param *MAOParam) ValidateApp(val string) (status bool, msg string) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
	}
	bodyStr := "{\"tool\":\"panosece/skantest\",\"artefact\":" + val + "}"

	client := &http.Client{Transport: tr}
	urlString := ""
	if strings.HasPrefix(param.MAOHost, "http") {
		urlString = param.MAOHost + "/inmem"
	} else {
		urlString = "http://" + param.MAOHost + "/inmem"
	}
	req, _ := http.NewRequest("POST", urlString, bytes.NewBuffer([]byte(bodyStr)))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", param.APIKey)

	res, err := client.Do(req)

	if err != nil {
		println("error opening MAP api server connection: " + err.Error())
		status = false
		msg = "error connecting to MAO server"
		return
	}
	defer res.Body.Close()

	if err != nil {
		status = false
		msg = err.Error()
		return
	} else {
		receivedBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			status = false
			msg = err.Error()
			return
		} else {
			println("received status code: " + res.Status)
			if strings.Contains(res.Status, "404 Not Found") {
				status = false
				msg = string(receivedBody)
				return
			}
			status = true
			msg = string(receivedBody)
		}
	}
	return
}
