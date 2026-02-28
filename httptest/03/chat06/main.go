package main

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"log"
	"os"
	"path/filepath"
)

// 基于 HTTPS的 证书扎钉 、 mTLS同

func main() {

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

	poolCA := loadCert("rootCA.key", "rootCA.cer")
	serverCert := loadCert("server.cer", "server.key")

	cfg := &tls.Config{
		Certificates: []tls.Certificate{*serverCert},
		VerifyConnection: func(state tls.ConnectionState) error {

			return nil
		},
		RootCAs:    poolCA,
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  nil,
		MinVersion: tls.VersionTLS12,
	}

}

/*

package main

import (
    "crypto/sha256"
    "crypto/tls"
    "crypto/x509"
    "encoding/base64"
    "fmt"
    "io"
    "net/http"
    "time"
)

// 预计算的公钥哈希（SPKI SHA-256，Base64 编码）
const pinnedPublicKeyHash = "base64_encoded_spki_hash_here"

func verifyPinnedPublicKey(state *tls.ConnectionState) error {
    if len(state.PeerCertificates) == 0 {
        return fmt.Errorf("no peer certificates")
    }

    cert := state.PeerCertificates[0]
    spkiHash := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
    hashBase64 := base64.StdEncoding.EncodeToString(spkiHash[:])

    if hashBase64 != pinnedPublicKeyHash {
        return fmt.Errorf("public key pinning failed: expected %s, got %s", pinnedPublicKeyHash, hashBase64)
    }

    return nil
}

func main() {
    tlsConfig := &tls.Config{
        VerifyConnection: verifyPinnedPublicKey,
        MinVersion:       tls.VersionTLS12,
    }

    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsConfig,
        },
        Timeout: 10 * time.Second,
    }

    resp, err := client.Get("https://example.com")
    if err != nil {
        fmt.Println("Request failed:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Println("Response:", string(body))
}

*/
