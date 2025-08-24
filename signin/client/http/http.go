package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var errCode = errors.New("resp status code error! ")

var _client *http.Client

func New() {
	_client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   5,
			MaxConnsPerHost:       10,
			IdleConnTimeout:       300 * time.Second,
			ResponseHeaderTimeout: 120 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Second * 30,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func Do(ctx context.Context, url string, meta *M) (*Resp, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(meta.Body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	if len(meta.Head) != 0 {
		for _, value := range meta.Head {
			req.Header.Set(value.Key, value.Value)
		}
	}

	resp, err := _client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errCode
	}
	defer resp.Body.Close()

	res := new(Resp)
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(res); err != nil {
		return nil, err
	}
	return res, nil
}

type Resp struct {
	Code int    `json:"code"`
	Data []byte `json:"data"`
}
