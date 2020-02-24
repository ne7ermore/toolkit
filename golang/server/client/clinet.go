package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var (
	cliOnce    sync.Once
	httpClient *http.Client
)

func getClient(disKA, disCompression bool, timeout int) *http.Client {
	cliOnce.Do(func() {
		hc := new(http.Client)
		hc.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DisableKeepAlives:   disKA,
			DisableCompression:  disCompression,
			TLSHandshakeTimeout: time.Duration(timeout) * time.Millisecond,
		}
		httpClient = hc
	})
	return httpClient
}

func NewPostJsonReq(url string, params map[string]interface{}) (*http.Request, error) {
	bytesData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return req, nil
}

func NewPostFormReq(url string, form url.Values) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

type client struct {
	client                *http.Client
	disKA, disCompression bool
	timeout               int
}

func InitClient(disKA, disCompression bool, timeout int) *client {
	c := new(client)
	c.client = getClient(disKA, disCompression, timeout)
	return c
}

func (c *client) Run(q *http.Request) (string, error) {
	resp, err := c.client.Do(q)
	if err != nil {
		return "", err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	str := (*string)(unsafe.Pointer(&respBytes))

	return *str, nil
}
