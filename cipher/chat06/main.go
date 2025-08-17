package main

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	bcryptTest()
	fmt.Println()
	argon2idTest()

}

func bcryptTest() {

	// CPU密集型

	password := "password"

	// 1. 使用 bcrypt.GenerateFromPassword 生成哈希。
	// 第二个参数是成本因子 (cost)。推荐值通常在 10 到 14 之间。
	// 成本每增加 1，计算时间大约翻一倍。

	cost := bcrypt.DefaultCost
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		panic(err)
	}

	fmt.Println("real passwd => ", password)
	fmt.Println("write db passwd => ", string(hashedPasswd))

	loginPwd := "password"

	if err = bcrypt.CompareHashAndPassword(hashedPasswd, []byte(loginPwd)); err != nil {
		panic(err)
	} else {
		fmt.Println("verify Ok!")
	}

	//bcrypt.GenerateFromPassword: 这个函数完成了所有繁重的工作——生成随机盐，并根据指定的成本进行哈希。
	//返回的字节串是一个自包含的格式，形如 $2a$10$N9qo8uLOickgx2ZMRZoMye.IKbeTj7jmwHW_sLwUnFDG./nI9I5pS，
	//其中包含了算法版本、成本和盐值。bcrypt.CompareHashAndPassword: 这是验证密码的唯一正确方法。
	//你永远不需要自己去解析哈希串中的盐。这个函数为你处理好了一切，并且内部使用了常量时间比较来防止时序攻击。

}

func scryptTest() {

	// scrypt CPU内存密集型

}

func argon2idTest() {

	// bcrypt.GenerateFromPassword: 这个函数完成了所有繁重的工作——生成随机盐，
	//并根据指定的成本进行哈希。返回的字节串是一个自包含的格式，
	//形如 $2a$10$N9qo8uLOickgx2ZMRZoMye.IKbeTj7jmwHW_sLwUnFDG./nI9I5pS，
	//其中包含了算法版本、成本和盐值。bcrypt.CompareHashAndPassword: 这是验证密码的唯一正确方法。
	//你永远不需要自己去解析哈希串中的盐。这个函数为你处理好了一切，并且内部使用了常量时间比较来防止时序攻击。

	password := "password"

	params := &ArgonParams{
		Memory:      64 * 1024, // 64MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	hashPasswd, err := generateHashPassword(password, params)
	if err != nil {
		panic(err)
	}

	fmt.Println("hash passwd => ", hashPasswd)

	ok, err := comparePasswdAndHash(password, hashPasswd)
	if err != nil {
		panic(err)
	}

	if !ok {
		panic("passwd and hash do not match")
	} else {
		fmt.Println("verify Ok!")
	}
}

type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func generateHashPassword(password string, argonParams *ArgonParams) (string, error) {
	salt := make([]byte, argonParams.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt,
		argonParams.Iterations,
		argonParams.Memory,
		argonParams.Parallelism,
		argonParams.KeyLength)

	b64salt := base64.StdEncoding.EncodeToString(salt)
	b64hash := base64.StdEncoding.EncodeToString(hash)

	var encodeHash strings.Builder
	encodeHash.WriteString("$argon2id$")
	encodeHash.WriteString(fmt.Sprintf("version=%d$", argon2.Version))
	encodeHash.WriteString(fmt.Sprintf("memory=%d,", argonParams.Memory))
	encodeHash.WriteString(fmt.Sprintf("iterations=%d,", argonParams.Iterations))
	encodeHash.WriteString(fmt.Sprintf("parallelism=%d$%s$%s", argonParams.Parallelism, b64salt, b64hash))
	return encodeHash.String(), nil
}

func comparePasswdAndHash(passwd, encodeHash string) (bool, error) {

	parts := strings.Split(encodeHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid encoded hash: %s", encodeHash)
	}

	// 验证版本号
	var version int
	_, err := fmt.Sscanf(parts[2], "version=%d", &version)
	if err != nil {
		return false, err
	}

	if version != argon2.Version {
		return false, fmt.Errorf("version failed: %s", encodeHash)
	}

	// 解析参数
	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "memory=%d,iterations=%d,parallelism=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, err
	}

	// 解析salt
	salt, err := base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	// 解析哈希
	decodeHash, err := base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	keyLength := uint32(len(decodeHash))
	// 使用输入的密码和提取出来的参数，重新计算哈希
	compareHash := argon2.IDKey([]byte(passwd), salt, iterations, memory, parallelism, keyLength)
	return hmac.Equal(decodeHash, compareHash), nil
}
