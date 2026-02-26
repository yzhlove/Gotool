package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// 接续 pem格式证书

func main() {

	data, err := os.ReadFile("/Users/yostar/Develop/Go/GoPath/src/rain.com/Gotool/httptest/03/certificate/rootCA.cer")
	if err != nil {
		panic(err)
	}

	block, reset := pem.Decode(data)
	fmt.Println(string(reset))

	fmt.Println("Type:", block.Type, block.Headers)

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(cert.AuthorityKeyId))
	fmt.Println(cert.IsCA)
	fmt.Println(cert.Subject)
	fmt.Println(cert.Issuer)

}
