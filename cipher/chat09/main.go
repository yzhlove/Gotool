package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func main() {

	x, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	pub := x.PublicKey().Bytes()
	fmt.Println("privKey ==> ", hex.EncodeToString(x.Bytes()))
	fmt.Println("pubKey ==> ", hex.EncodeToString(pub))

	x2, err := ecdh.X25519().NewPrivateKey(x.Bytes())
	if err != nil {
		panic(err)
	}

	pub2 := hex.EncodeToString(x2.PublicKey().Bytes())
	fmt.Println("pubKey2 ==> ", pub2)
	fmt.Println("privKey2 ==> ", hex.EncodeToString(x2.Bytes()))

}
