package api

import (
	"fmt"
	"github.com/yzhlove/Gotool/bing/module/log"
	"net/url"
	"testing"
)

func Test_Wallpaper(t *testing.T) {

	log.New()

	_, err := GetWallpaper()
	if err != nil {
		t.Fatal(err)
		return
	}

}

func Test_2(t *testing.T) {

	s := "/th?id=OHR.PandaSnow_ZH-CN5981854301"

	ret, err := url.Parse(s)
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(ret.Scheme, ret.Host, ret.Path, ret.Query())

}
