package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
)

// HMAC 指纹

func main() {

	// 未被篡改的消息和MAC信息
	a, b := mockClient()
	mockServer(a, b)

	// 篡改消息
	a, b = mockClient()
	a = "this is attack message!"
	mockServer(a, b)

	// 篡改MAC信息
	a, b = mockClient()
	m := md5.New()
	io.Copy(m, bytes.NewBuffer([]byte(a)))
	b = string(m.Sum(nil))
	mockServer(a, b)

	// 篡改消息和MAC信息
	a, b = mockClient()
	m = md5.New()
	io.Copy(m, bytes.NewBuffer([]byte("hello world!!!")))
	b = string(m.Sum(nil))
	mockServer(a, b)

}

func generateMAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func verifyMAC(message, messageMAC, key []byte) bool {

	newMessageMAC := generateMAC(message, key)
	// 使用 hmac.Equal() 函数比较两个MAC值是否相等，防止时序攻击
	// 时序攻击是指攻击者利用计算时间差异进行攻击的一种手段。
	// 例如，攻击者可以通过测量加密算法的运行时间，来猜测明文的内容。
	// 通过比较两个MAC值是否相等，可以避免时序攻击。
	// 什么是时序攻击 (Timing Attack)？
	//这是一种旁道攻击。
	//如果一个比较函数在发现第一个不匹配的字节时就立即返回（“短路”），
	//那么比较两个字符串所花费的时间就会因它们有多少位前缀是相同的而不同。
	//攻击者可以通过精确测量这些微小的响应时间差异，逐个字节地猜出整个秘密值。
	//hmac.Equal 通过确保无论内容是否匹配，比较过程都耗费相同的时间（常量时间），从而有效抵御此类攻击。
	return hmac.Equal(messageMAC, newMessageMAC)
}

var secret = []byte("*#06#*")

func mockClient() (string, string) {

	msg := []byte(`msg="hello+world&user=123456&time=1234567890"`)
	macMsg := generateMAC(msg, secret)
	return string(msg), string(macMsg)
}

func mockServer(msg, macMsg string) {
	if verifyMAC([]byte(msg), []byte(macMsg), secret) {
		fmt.Println("msg is OK!")
	} else {
		fmt.Println("msg is invalid!")
	}
}
