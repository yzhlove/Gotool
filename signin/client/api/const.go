package api

func trim(str string) string {
	if len(str) > 0 {
		if str[len(str)-1] == '/' {
			str = str[:len(str)-1]
		}
	}
	return str
}

func getIkeApi(url string) string {
	return trim(url) + "/ike"
}
