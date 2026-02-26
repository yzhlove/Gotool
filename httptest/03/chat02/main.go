package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

// 创建证书

func main() {

	prviateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		log.Fatal(err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{"CN"},
			StreetAddress:      []string{"Shanghai"},
			Locality:           []string{"Shanghai"},
			Organization:       []string{"MyTestOrg"},
			OrganizationalUnit: []string{"MyTestOrg"},
			CommonName:         "www.yzhdomain.com",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &prviateKey.PublicKey, prviateKey)
	if err != nil {
		log.Fatal(err)
	}

	block := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	if err = os.WriteFile("/Users/yostar/Develop/Go/GoPath/src/rain.com/Gotool/httptest/03/chat02/server.cer", block, 0644); err != nil {
		log.Fatal(err)
	}

}
