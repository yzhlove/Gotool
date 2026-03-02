package main

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// 基于 HTTPS的 证书扎钉 、 mTLS同

func main() {
	go httpsServer()
	time.Sleep(time.Second * 2)
	httpClientOk()
	time.Sleep(time.Second * 2)
}

var (
	rootPATH = "/Users/yostar/Develop/Go/GoPath/src/rain.com/Gotool/httptest/03/chat04"
)

func loadCACert(cert string) *x509.CertPool {

	pemCA, err := os.ReadFile(filepath.Join(rootPATH, cert))
	if err != nil {
		log.Fatal(err)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCA)
	return pool
}

func loadCert(keyPath, certPath string) *tls.Certificate {

	pemKey, err := os.ReadFile(filepath.Join(rootPATH, keyPath))
	if err != nil {
		log.Fatal(err)
	}

	pemCert, err := os.ReadFile(filepath.Join(rootPATH, certPath))
	if err != nil {
		log.Fatal(err)
	}

	certificate, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		log.Fatal(err)
	}
	return &certificate
}
func mockGetCertHash(cert string) string {

	pemData, err := os.ReadFile(filepath.Join(rootPATH, cert))
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(pemData)
	if block.Type != "CERTIFICATE" {
		log.Fatal("invalid cert")
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	// 计算证书的哈希，通常这个步骤不会出现在代理里面，都是提前计算好的
	hash := sha256.Sum256(certificate.RawSubjectPublicKeyInfo)
	return hex.EncodeToString(hash[:])
}

func httpsServer() {

	poolCA := loadCACert("rootCA.cer")
	serverCert := loadCert("server.key", "server.cer")

	// 服务器验证客户端证书，正常情况下，服务器不做这个操作，证书扎钉只在客户端检查就行。本次为演示流程
	cliCertHash := mockGetCertHash("client.cer")

	cfg := &tls.Config{
		Certificates: []tls.Certificate{*serverCert},
		VerifyConnection: func(state tls.ConnectionState) error {
			if len(state.PeerCertificates) == 0 {
				return fmt.Errorf("no peer certificates")
			}

			cert := state.PeerCertificates[0]
			fmt.Println("client ---> 1 ", cert.Version)
			fmt.Println("client ---> 2 ", cert.Subject)
			fmt.Println("client ---> 3 ", cert.ExtKeyUsage)
			fmt.Println("client ---> 4 ", cert.IsCA)

			sha := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
			if cliCertHash != hex.EncodeToString(sha[:]) {
				return fmt.Errorf("server invalid client certificate")
			}
			return nil
		},
		ClientAuth: tls.RequireAndVerifyClientCert, // mTLS过程中，如果强制要求客户端传证书，则需要设置ClientAuth为tls.RequireAndVerifyClientCert 以及 ClientCAs
		ClientCAs:  poolCA,
		MinVersion: tls.VersionTLS12,
	}

	h := http.NewServeMux()
	h.HandleFunc("/", handler)
	s := &http.Server{
		Addr:      ":443",
		TLSConfig: cfg,
		Handler:   h,
	}
	fmt.Println("https listen to ", s.Addr)
	if err := s.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Strict-Transport-Security", "max-age=300; includeSubDomains;")

	resp := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		0,
		fmt.Sprintf("Ok! %s", time.Now().Format(time.RFC3339)),
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
	}
}

func httpClientOk() {

	poolCA := loadCACert("rootCA.cer")
	clientCert := loadCert("client.key", "client.cer")

	serverHash := mockGetCertHash("server.cer")

	// serverHash
	// serverHash = "123456789@hash"

	tlsCfg := &tls.Config{
		ServerName:   "www.yzhdomain.com",
		RootCAs:      poolCA, // 当 rootCAs为空的时候，客户端通常会加载系统证书池
		Certificates: []tls.Certificate{*clientCert},
		MinVersion:   tls.VersionTLS12,
		VerifyConnection: func(state tls.ConnectionState) error {
			if len(state.PeerCertificates) == 0 {
				return fmt.Errorf("no peer certificates")
			}

			cert := state.PeerCertificates[0]
			hash := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
			if serverHash != hex.EncodeToString(hash[:]) {
				return fmt.Errorf("client invalid server certificate")
			}
			return nil
		},
	}

	cc := http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
	}

	resp, err := cc.Get("https://www.yzhdomain.com")
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(data))
}
