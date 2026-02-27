package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// 使用 GO 生成 CA mTLS 所需的证书

func main() {

	//generateCert()
	checkCert()
}

func generateCert() {
	caKey, caCert := genRootCA()
	genServerCert(caKey, caCert)
	genClientCert(caKey, caCert)
}

type bytesWrap []byte

func (b bytesWrap) Save(file, blockType string) error {
	value := pem.EncodeToMemory(&pem.Block{Type: blockType, Bytes: b})
	return os.WriteFile(filepath.Join(rootPATH, file), value, os.ModePerm)
}

var (
	pkixName = pkix.Name{
		Country:            []string{"CN"},
		Province:           []string{"Shanghai"},
		Locality:           []string{"Shanghai"},
		Organization:       []string{"MyTestOrg"},
		OrganizationalUnit: []string{"MyTestUnit"},
		CommonName:         "www.yzhdomain.com",
	}
	serverDNSs = []string{"www.yzhdomain.com", "localhost"}
	serverIPs  = []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("10.155.120.47"),
		net.ParseIP("10.155.90.131"),
	}
	rootPATH = "/Users/yostar/Develop/Go/GoPath/src/rain.com/Gotool/httptest/03/chat04"
)

func genRootCA() (*ecdsa.PrivateKey, *x509.Certificate) {

	// CA 私钥
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// CA 证书编号
	setialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		log.Fatal(err)
	}

	temp := &x509.Certificate{
		SerialNumber:          setialNumber,
		Subject:               pkixName,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true, // 这是关键
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCRLSign | x509.KeyUsageCertSign, // 这是关键 ， CA 证书必须具备 签名和 CA 签名
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
	}

	// 自签证书的 parent 和 signer 都是自身
	derBytes, err := x509.CreateCertificate(rand.Reader, temp, temp, &caKey.PublicKey, caKey)
	if err != nil {
		log.Fatal(err)
	}

	caCer, err := x509.ParseCertificate(derBytes)
	if err != nil {
		log.Fatal(err)
	}

	// 保存密钥
	caKeyBytes, err := x509.MarshalECPrivateKey(caKey)
	if err != nil {
		log.Fatal(err)
	}

	err = bytesWrap(caKeyBytes).Save("rootCA.key", "EC PRIVATE KEY")
	if err != nil {
		log.Fatal(err)
	}

	// 保存证书
	err = bytesWrap(derBytes).Save("rootCA.cer", "CERTIFICATE")
	if err != nil {
		log.Fatal(err)
	}
	return caKey, caCer
}

func genServerCert(caKey *ecdsa.PrivateKey, caCer *x509.Certificate) {

	sevKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		log.Fatal(err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkixName,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		DNSNames:     serverDNSs,
		IPAddresses:  serverIPs,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	// 生成服务器证书
	sevCertBytes, err := x509.CreateCertificate(rand.Reader, tmpl, caCer, &sevKey.PublicKey, caKey)
	if err != nil {
		log.Fatal(err)
	}

	// 保存私钥
	sevBytes, err := x509.MarshalECPrivateKey(sevKey)
	if err != nil {
		log.Fatal(err)
	}

	err = bytesWrap(sevBytes).Save("server.key", "EC PRIVATE KEY")
	if err != nil {
		log.Fatal(err)
	}

	// 保存证书
	err = bytesWrap(sevCertBytes).Save("server.cer", "CERTIFICATE")
	if err != nil {
		log.Fatal(err)
	}
}

func genClientCert(caKey *ecdsa.PrivateKey, caCer *x509.Certificate) {

	cliKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		log.Fatal(err)
	}

	pkixNamecp := pkixName
	pkixNamecp.CommonName = "rain@12138.com"

	tmpl := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkixNamecp,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		DNSNames:     serverDNSs,
		IPAddresses:  serverIPs,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	// 生成客户端证书
	cliCertBytes, err := x509.CreateCertificate(rand.Reader, tmpl, caCer, &cliKey.PublicKey, caKey)
	if err != nil {
		log.Fatal(err)
	}

	err = bytesWrap(cliCertBytes).Save("client.cer", "CERTIFICATE")
	if err != nil {
		log.Fatal(err)
	}

	cliBytes, err := x509.MarshalECPrivateKey(cliKey)
	if err != nil {
		log.Fatal(err)
	}

	err = bytesWrap(cliBytes).Save("client.key", "EC PRIVATE KEY")
	if err != nil {
		log.Fatal(err)
	}
}

func verifyCertByRoot(root *x509.Certificate, cert *x509.Certificate, opts x509.VerifyOptions) error {
	p := x509.NewCertPool()
	p.AddCert(root)
	opts.Roots = p

	_, err := cert.Verify(opts)
	return err
}

func loadCert(keyPath, certPath string) (*ecdsa.PrivateKey, *x509.Certificate) {

	certBytes, err := os.ReadFile(filepath.Join(rootPATH, certPath))
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(certBytes)
	if block.Type != "CERTIFICATE" {
		log.Fatal("key type error")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	privateBytes, err := os.ReadFile(filepath.Join(rootPATH, keyPath))
	if err != nil {
		log.Fatal(err)
	}

	keyBlock, _ := pem.Decode(privateBytes)
	if keyBlock.Type != "EC PRIVATE KEY" {
		log.Fatal("key type error")
	}

	key, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	return key, cert
}

func checkCert() {

	rootCAKey, rootCACer := loadCert("rootCA.key", "rootCA.cer")
	serverKey, serverCer := loadCert("server.key", "server.cer")
	clientKey, clientCer := loadCert("client.key", "client.cer")

	// 验证域名
	sOpts := x509.VerifyOptions{
		CurrentTime: time.Now(),
		KeyUsages:   []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSName:     "www.yzhdomain.com",
	}

	if err := verifyCertByRoot(rootCACer, serverCer, sOpts); err != nil {
		panic(fmt.Sprintf("check error: %+v", err))
	}

	cOPts := x509.VerifyOptions{
		CurrentTime: time.Now(),
		KeyUsages:   []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	if err := verifyCertByRoot(rootCACer, clientCer, cOPts); err != nil {
		panic(fmt.Sprintf("check error: %+v", err))
	}

	// 验证 IP
	ip := net.ParseIP("10.155.120.47")
	if err := serverCer.VerifyHostname(ip.String()); err != nil {
		panic(fmt.Sprintf("check error: %+v", err))
	}

	_ = rootCAKey
	_ = serverKey
	_ = clientKey
}
