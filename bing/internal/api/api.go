package api

import (
	"encoding/json"
	"fmt"
	"github.com/yzhlove/Gotool/bing/internal/config"
	"github.com/yzhlove/Gotool/bing/module/client"
	"github.com/yzhlove/Gotool/bing/module/log"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

func getWallpaperUrl() string {

	var u = new(url.URL)
	var values = url.Values{}
	// 请求参数：
	//   format：返回的数据文件格式，由json或者xml文档，我这里选择xml文档，更适合解析
	//   idx：相当于请求哪个图片，填0就是今天的，填1就是昨天的，以此类推
	//   n：请求图片的数量，1就是今天的，2就是今天和昨天的
	//   mkt：默认填zh-CN
	// 	 uhd：1	若为空或0，返回1080P结果，若为 1 则返回4K结果，==从2019.05.10开始支持==
	values.Set("format", "js")
	values.Set("idx", "0")
	values.Set("n", "1")
	values.Set("mkt", string(config.China))
	values.Set("uhd", "1")

	u.Scheme = config.Scheme
	u.Host = config.Host
	u.Path = config.Path
	u.RawQuery = values.Encode()

	return u.String()
}

func downWallpaperUrl(address string) (string, error) {

	ret, err := url.Parse(address)
	if err != nil {
		return "", err
	}

	str := ret.Query().Get("id")
	if len(str) == 0 {
		return "", fmt.Errorf("url:%s not found key['id']! ", address)
	}

	var u = &url.URL{}
	u.Scheme = config.Scheme
	u.Host = config.Host
	u.Path = ret.Path

	var values = url.Values{}
	values.Set("id", fmt.Sprintf("%s_UHD.jpg", str))
	values.Set("rf", "LaDigue_UHD.jpg")
	u.RawQuery = values.Encode()

	return u.String(), nil
}

func GetWallpaper() (*Msg, error) {

	addressUrl := getWallpaperUrl()

	log.Debug("address url", slog.String("address_url", addressUrl))

	req, err := http.NewRequest(http.MethodGet, addressUrl, nil)
	if err != nil {
		return nil, err
	}

	var jsonData []byte
	err = client.HttpGet(req, func(reader io.Reader) error {
		if jsonData, err = io.ReadAll(reader); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var resp = &BingResp{}
	if err = json.Unmarshal(jsonData, resp); err != nil {
		return nil, err
	}

	if len(resp.Images) == 0 {
		return nil, fmt.Errorf("images is empty! ")
	}

	downUrl, err := downWallpaperUrl(resp.Images[0].UrlBase)
	if err != nil {
		return nil, err
	}

	log.Debug("download url", slog.String("download_url", downUrl))

	req, err = http.NewRequest(http.MethodGet, downUrl, nil)
	if err != nil {
		return nil, err
	}

	msg := &Msg{
		Title: resp.Images[0].Title,
		Auth:  resp.Images[0].Auth,
	}

	err = client.HttpGet(req, func(reader io.Reader) error {
		if msg.Content, err = io.ReadAll(reader); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return msg, nil
}
