package config

const (
	Scheme = "https"
	Host   = "cn.bing.com"
	Path   = "HPImageArchive.aspx"
)

type Country string

const (
	China   Country = "zh-CN" // 中国
	Japan   Country = "ja-JP" // 日本
	Germany Country = "de-DE" // 德国
	Canada  Country = "en-CA" // 加拿大
	England Country = "en-GB" // 英国
	India   Country = "en-IN" // 印度
	USA     Country = "en-US" // 美国
	France  Country = "fr-FR" // 法国
	Italy   Country = "it-IT" // 意大利
)
