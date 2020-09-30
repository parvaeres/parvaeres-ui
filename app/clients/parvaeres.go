package clients

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"go8s/app/models"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type ParvaeresParam struct {
	APIKey     string
	APIHost    string
	APIPort    string
	APIVersion string
}

func (param *ParvaeresParam) GetApplication(id string) (status bool, msg string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
	}

	client := &http.Client{Transport: tr}
	urlString := ""
	if strings.HasPrefix(param.APIHost, "http") {
		if len(param.APIPort) > 0 {
			urlString = param.APIHost + ":" + param.APIPort + "/" + param.APIVersion + "/deployment/" + id
		} else {
			urlString = param.APIHost + "/" + param.APIVersion + "/deployment/" + id
		}
	} else {
		if len(param.APIPort) > 0 {
			urlString = "http://" + param.APIHost + ":" + param.APIPort + "/" + param.APIVersion + "/deployment/" + id
		} else {
			urlString = "http://" + param.APIHost + "/" + param.APIVersion + "/deployment/" + id
		}
	}
	req, _ := http.NewRequest("GET", urlString, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", param.APIKey)

	res, err := client.Do(req)
	if err != nil {
		println("error opening parvaeres api server connection: " + err.Error())
		status = false
		msg = "error connecting to api server"
		return
	}
	defer res.Body.Close()

	receivedBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		status = false
		msg = err.Error()
		return
	} else {
		status = true
		msg = string(receivedBody)
	}
	return
}

func (param *ParvaeresParam) RegisterApplication(val models.ParvaeresApplicationData) (status bool, msg string) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
	}
	body, _ := json.Marshal(val)

	client := &http.Client{Transport: tr}
	urlString := ""
	if strings.HasPrefix(param.APIHost, "http") {
		if len(param.APIPort) > 0 {
			urlString = param.APIHost + ":" + param.APIPort + "/" + param.APIVersion + "/deployment"
		} else {
			urlString = param.APIHost + "/" + param.APIVersion + "/deployment"
		}
	} else {
		if len(param.APIPort) > 0 {
			urlString = "http://" + param.APIHost + ":" + param.APIPort + "/" + param.APIVersion + "/deployment"
		} else {
			urlString = "http://" + param.APIHost + "/" + param.APIVersion + "/deployment"
		}
	}
	req, _ := http.NewRequest("POST", urlString, bytes.NewBuffer(body))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", param.APIKey)

	res, err := client.Do(req)
	if err != nil {
		println("error opening parvaeres api server connection: " + err.Error())
		status = false
		msg = "error connecting to api server"
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
			status = true
			msg = string(receivedBody)
		}
	}
	return
}
