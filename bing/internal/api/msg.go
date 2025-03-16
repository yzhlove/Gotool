package api

type BingResp struct {
	Images []struct {
		Auth    string `json:"copyright"`
		UrlBase string `json:"urlbase"`
		Title   string `json:"title"`
	} `json:"images"`
}

type Msg struct {
	Content []byte
	Title   string
	Auth    string
}
