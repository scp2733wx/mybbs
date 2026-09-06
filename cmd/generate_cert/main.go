// generate_cert 生成用于本地 HTTPS 测试的自签名证书。
// 用法: go run ./cmd/generate_cert
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("生成私钥失败: %v", err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		log.Fatalf("生成序列号失败: %v", err)
	}

	tmpl := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   "localhost",
			Organization: []string{"mybbs dev"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	der, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("创建证书失败: %v", err)
	}

	if err := writePEM("cert.pem", "CERTIFICATE", der); err != nil {
		log.Fatalf("写入 cert.pem 失败: %v", err)
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("编码私钥失败: %v", err)
	}
	if err := writePEM("key.pem", "PRIVATE KEY", keyBytes); err != nil {
		log.Fatalf("写入 key.pem 失败: %v", err)
	}

	fmt.Println("已生成自签名证书:")
	fmt.Println("  cert.pem (证书)")
	fmt.Println("  key.pem  (私钥)")
	fmt.Println("浏览器首次访问会警告，请手动信任该证书。")
}

func writePEM(path, blockType string, bytes []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return pem.Encode(f, &pem.Block{Type: blockType, Bytes: bytes})
}
