package handler

type Msg struct {
	Code int    `json:"code,omitempty"`
	Data []byte `json:"data,omitempty"`
}

func buildErr(err error) Msg {
	return Msg{
		Code: -1,
		Data: []byte(err.Error()),
	}
}

func buildOk(data []byte) Msg {
	return Msg{
		Code: 0,
		Data: data,
	}
}
