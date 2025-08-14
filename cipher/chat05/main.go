package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/hkdf"
	"io"
)

func main() {
	ecdsaTest()
	hKDFTest()
}

func ecdsaTest() {

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	message := "hello world"
	hash := sha256.Sum256([]byte(message))

	sg, err := ecdsa.SignASN1(rand.Reader, privKey, hash[:])
	if err != nil {
		panic(err)
	}

	if !ecdsa.VerifyASN1(&privKey.PublicKey, hash[:], sg) {
		panic("check error")
	} else {
		fmt.Println("check Ok! ")
	}

}

func hKDFTest() {

	secret := []byte("aaaa")
	slat := []byte("bbb")
	info := []byte("what are you doing? ")
	reader := hkdf.New(sha256.New, secret, slat, info)

	bytes := make([]byte, 32)
	if _, err := io.ReadFull(reader, bytes); err != nil {
		panic(err)
	}
	fmt.Printf("secrt = > %x\n", bytes)

}
