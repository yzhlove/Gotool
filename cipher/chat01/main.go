package main

import "fmt"

// 01 简单的XOR加密

func main() {
	//test1()
	test2()
}

func xorCipher(key, text []byte) []byte {
	output := make([]byte, len(text))
	for i, v := range text {
		output[i] = v ^ key[i%len(key)]
	}
	return output
}

func test1() {
	key := []byte("this is secret!")
	text := []byte("hello world")
	encrypt := xorCipher(key, text)
	fmt.Println(string(encrypt))
	decrypt := xorCipher(key, encrypt)
	fmt.Println(string(decrypt))
}

func test2() {
	// 密码重用的灾难

	key := []byte("this is secret!")
	text1 := []byte("hello world!")
	text2 := []byte("hello gopher!")

	fmt.Println("text1 string => ", string(text1))
	fmt.Println("text2 string => ", string(text2))

	// 用同一个秘钥加密的两条消息
	ciphere1 := xorCipher(key, text1)
	ciphere2 := xorCipher(key, text2)

	fmt.Println("cipher1 string => ", ciphere1)
	fmt.Println("cipher2 string => ", ciphere2)

	// 攻击者视角
	minLength := min(len(ciphere1), len(ciphere2))
	attCipherText := xorCipher(ciphere1[:minLength], ciphere2[:minLength])

	// 攻击者可以很容易的获得到两端明文的异或结果，自然也很容易推断出原始明文内容
	fmt.Println("attacker's view => ", attCipherText)

	// checker
	minLength = min(len(text1), len(text2))
	text := xorCipher(text1[:minLength], text2[:minLength])

	fmt.Println("checker's view => ", text)

}
