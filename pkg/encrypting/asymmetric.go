package encrypting

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

const keyBytesLength = 4096

// MakePrivateKey Create private key.
func MakePrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, keyBytesLength)
}

func makeCertTemplate() *x509.Certificate {
	return &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658), //nolint:gomnd,mnd
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback}, //nolint:gomnd,mnd
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0), //nolint:gomnd,mnd
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}
}

// MakeCert Create certificate.
func MakeCert(privKey *rsa.PrivateKey) ([]byte, error) {
	cert := makeCertTemplate()

	return x509.CreateCertificate(rand.Reader, cert, cert, &privKey.PublicKey, privKey)
}

// GenerateCertificate Generate certificate.
func GenerateCertificate(host string) (*tls.Certificate, error) {
	privateKey, err := MakePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	tpl := makeCertTemplate()
	certBytes, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certFile, err := os.Create("certs/server.crt")
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate file: %w", err)
	}
	defer certFile.Close()
	certPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}
	if err := pem.Encode(certFile, certPEMBlock); err != nil {
		return nil, fmt.Errorf("failed to encode certificate: %w", err)
	}

	keyFile, err := os.Create("certs/server.key")
	if err != nil {
		return nil, fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyFile.Close()
	keyPEMBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(keyFile, keyPEMBlock); err != nil {
		return nil, fmt.Errorf("failed to encode key: %w", err)
	}

	cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate and key: %w", err)
	}

	return &cert, nil
}
