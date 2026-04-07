package main

import (
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"sync"
)

type rwPair struct {
	io.Reader
	io.Writer
}

func main() {

	sp, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	LongTremData = sp.PublicKey().Bytes()
	_sevPrivateKey = sp.Bytes()

	// client→server: pr1/pw1, server→client: pr2/pw2
	pr1, pw1 := io.Pipe()
	pr2, pw2 := io.Pipe()
	serverRW := rwPair{pr1, pw2}
	clientRW := rwPair{pr2, pw1}
	var wg sync.WaitGroup
	wg.Go(func() {
		server(serverRW)
		pw2.Close()
	})
	wg.Go(func() {
		client(clientRW)
		pw1.Close()
	})
	wg.Wait()

}

var LongTremData []byte
var _sevPrivateKey []byte
var InfoMessage = "test-info-message"

func server(rw io.ReadWriter) {

	p1 := packet{}

	fmt.Println("s r p1")

	if err := p1.read(rw); err != nil {
		panic(err)
	}

	p2 := packet{}

	fmt.Println("s r p2")
	if err := p2.read(rw); err != nil {
		panic(err)
	}

	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	p3 := packet{
		data: priv.PublicKey().Bytes(),
	}

	fmt.Println("s w p3")

	if err := p3.write(rw); err != nil {
		panic(err)
	}

	fmt.Println("s compute")

	cliPub, err := ecdh.X25519().NewPublicKey(p1.data)
	if err != nil {
		panic(err)
	}

	s1, err := priv.ECDH(cliPub)
	if err != nil {
		panic(err)
	}

	spriv2, err := ecdh.X25519().NewPrivateKey(_sevPrivateKey)
	if err != nil {
		panic(err)
	}

	s2, err := spriv2.ECDH(cliPub)
	if err != nil {
		panic(err)
	}

	// derive secret key
	secret := make([]byte, len(s1)+len(s2))
	copy(secret, s1)
	copy(secret[len(s1):], s2)

	secretKey, err := hkdf.Key(sha256.New, secret, p2.data, InfoMessage, 32)
	if err != nil {
		panic(err)
	}

	fmt.Println("server secret key => ", hex.EncodeToString(secretKey))
}

func client(rw io.ReadWriter) {

	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	// write public key to server
	p1 := packet{
		data: priv.PublicKey().Bytes(),
	}

	fmt.Println("c w p1")

	if err = p1.write(rw); err != nil {
		panic(err)
	}

	// write salt to server
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err)
	}

	p2 := packet{
		data: salt,
	}

	fmt.Println("c w p2")

	if err = p2.write(rw); err != nil {
		panic(err)
	}

	// read server public key
	p3 := packet{}

	fmt.Println("c r p3")

	if err := p3.read(rw); err != nil {
		panic(err)
	}

	fmt.Println("c compoute")

	sevPub, err := ecdh.X25519().NewPublicKey(p3.data)
	if err != nil {
		panic(err)
	}

	sevLongTremPub, err := ecdh.X25519().NewPublicKey(LongTremData)
	if err != nil {
		panic(err)
	}

	// compute shared secret
	s1, err := priv.ECDH(sevPub)
	if err != nil {
		panic(err)
	}

	// compute shared secret
	s2, err := priv.ECDH(sevLongTremPub)
	if err != nil {
		panic(err)
	}

	// derive secret key
	secret := make([]byte, len(s1)+len(s2))
	copy(secret, s1)
	copy(secret[len(s1):], s2)

	secretKey, err := hkdf.Key(sha256.New, secret, salt, InfoMessage, 32)
	if err != nil {
		panic(err)
	}

	fmt.Println("cli secret key => ", hex.EncodeToString(secretKey))
}

type packet struct {
	data []byte
}

func (p *packet) write(w io.Writer) error {
	l := len(p.data)
	if l > math.MaxUint16 {
		return fmt.Errorf("packet too large")
	}

	wdata := make([]byte, 2+l)
	binary.BigEndian.PutUint16(wdata[:2], uint16(l))
	copy(wdata[2:], p.data)
	_, err := w.Write(wdata)
	return err
}

func (p *packet) read(r io.Reader) error {
	dataLen := make([]byte, 2)
	if _, err := io.ReadFull(r, dataLen); err != nil {
		return err
	}
	l := binary.BigEndian.Uint16(dataLen)
	p.data = make([]byte, l)
	_, err := io.ReadFull(r, p.data)
	return err
}
