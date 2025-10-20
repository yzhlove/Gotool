package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/bytedance/sonic"
)

func toJson(path string) (map[string]any, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := make(map[string]any, 4)
	if err = sonic.ConfigDefault.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func getContent(value map[string]any) []string {
	var metas = make([]string, 64)
	if p1, ok := value["dat"]; ok {
		p2 := p1.(map[string]any)
		if p3, ok := p2["list"]; ok {
			if p4, ok := p3.([]any); ok {
				for _, pt := range p4 {
					if ret, ok := pt.(map[string]any); ok {
						if content, ok := ret["content"]; ok {
							metas = append(metas, content.(string))
						}
					}
				}
			}
		}
	}
	return metas
}

type Content struct {
	Timestamp int64          `json:"timestamp,omitempty"`
	Time      string         `json:"time,omitempty"`
	Level     string         `json:"level,omitempty"`
	Msg       string         `json:"msg,omitempty"`
	UID       uint64         `json:"uid,omitempty"`
	SdkUID    string         `json:"sdkuid,omitempty"`
	Nick      string         `json:"nick,omitempty"`
	Req       string         `json:"req,omitempty"`
	Stack     string         `json:"stack,omitempty"`
	Extra     map[string]any `json:"extra,omitempty"`
}

func (c *Content) optTime() {
	if tm, err := time.Parse(time.RFC3339, c.Time); err == nil {
		c.Timestamp = tm.Unix()
	}
}

func (c *Content) String() string {
	//var sb strings.Builder
	//sb.WriteString(fmt.Sprintf("Time:%s ", c.Time))
	//sb.WriteString(fmt.Sprintf("Level:%s ", c.Level))
	//sb.WriteString(fmt.Sprintf("Msg:%s ", c.Msg))
	//sb.WriteString(fmt.Sprintf("UID:%d ", c.UID))
	//sb.WriteString(fmt.Sprintf("SdkUID:%s ", c.SdkUID))
	//sb.WriteString(fmt.Sprintf("Nick:%s ", c.Nick))
	//sb.WriteString(fmt.Sprintf("Req:%s ", c.Req))
	//if c.Level == "ERROR" {
	//	sb.WriteString(fmt.Sprintf("Stack:%s ", c.Time))
	//}
	//sb.WriteString(goutil.String(c.Extra))
	//return sb.String()

	value, _ := sonic.MarshalString(c)
	return value
}

func toContent(value string) *Content {

	var data = make(map[string]any, 16)
	if err := sonic.UnmarshalString(value, &data); err != nil {
		slog.Error("to content failed! ", slog.Any("error", err))
		return nil
	}

	content := new(Content)
	if err := sonic.UnmarshalString(value, content); err != nil {
		slog.Error("to content failed! ", slog.Any("error", err))
		return nil
	}

	delete(data, "time")
	delete(data, "level")
	delete(data, "msg")
	delete(data, "uid")
	delete(data, "sdkuid")
	delete(data, "nick")
	delete(data, "stack")
	content.Extra = data
	return content
}
