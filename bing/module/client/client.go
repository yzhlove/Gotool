package client

import (
	"crypto/tls"
	"io"
	"net/http"
	"sync"
	"time"
)

var client *http.Client
var once sync.Once

func getHTTPClient() *http.Client {
	once.Do(func() {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				TLSHandshakeTimeout:   time.Second * 30,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   5,
				MaxConnsPerHost:       10,
				IdleConnTimeout:       time.Second * 10,
				ExpectContinueTimeout: time.Second * 5,
				ProxyConnectHeader:    nil,
				GetProxyConnectHeader: nil,
			},
			Timeout: time.Second * 30,
		}
	})
	return client
}

func HttpGet(req *http.Request, callback func(reader io.Reader) error) error {

	cc := getHTTPClient()
	resp, err := cc.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if callback != nil {
		if err = callback(resp.Body); err != nil {
			return err
		}
	}
	return nil
}
