package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"
)

// 创建一个用于简单tls验证的证书

func main() {
	//createCert()
	priv, cert := loadCert()
	verifyPrivAndCert(priv, cert)
	httpTLS(priv, cert)
}

var rootPath = "/Users/yostar/Develop/Go/GoPath/src/rain.com/Gotool/httptest/03/chat03"

func createCert() {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
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
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// 私钥
	privKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	privateKeyFile := rootPath + "/private.key"
	err = os.WriteFile(privateKeyFile, pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privKeyBytes,
	}), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// 证书
	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	certFile := rootPath + "/certificate.cer"
	err = os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func loadCert() (*ecdsa.PrivateKey, *x509.Certificate) {

	privData, err := os.ReadFile(rootPath + "/private.key")
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(privData)
	if block.Type != "EC PRIVATE KEY" {
		log.Fatal("invalid private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	certData, err := os.ReadFile(rootPath + "/certificate.cer")
	if err != nil {
		log.Fatal(err)
	}

	block, _ = pem.Decode(certData)
	if block.Type != "CERTIFICATE" {
		log.Fatal("invalid certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	return privateKey, cert
}

func verifyPrivAndCert(priv *ecdsa.PrivateKey, cert *x509.Certificate) {
	if !priv.PublicKey.Equal(cert.PublicKey) {
		log.Fatal("private key and certificate do not match")
	}
}

func httpTLS(priv *ecdsa.PrivateKey, cert *x509.Certificate) {

	keyPEMBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	if keyPEMBytes == nil {
		log.Fatal("certificate is nil")
	}

	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		log.Fatal(err)
	}

	privDataBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	})

	tlsCert, err := tls.X509KeyPair(keyPEMBytes, privDataBytes)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		MinVersion:   tls.VersionTLS12,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Strict-Transport-Security", "max-age=300; includeSubDomains;")

		x := struct {
			Time    time.Time `json:"time"`
			Message string    `json:"message"`
		}{
			Time:    time.Now(),
			Message: "https server run succeed! ",
		}

		e := json.NewEncoder(w)
		e.SetEscapeHTML(false)
		if err = e.Encode(x); err != nil {
			log.Fatal(err)
		}
	})

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
		Handler:   mux,
	}
	fmt.Println("start http server: ", server.Addr)
	if err = server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
