package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

type EtcdCert struct {
	subject pkix.Name
	config  CertConfig
}

type CertConfig struct {
	SubjectConfig    Subject
	Expire           time.Duration
	CaKeyFilePath    string
	CaCertFilePath   string
	EtcdKeyFilePath  string
	EtcdCertFilePath string
	EtcdCsrFilePath  string
	DnsList          []string
}

func NewEtcdCert(config CertConfig) *EtcdCert {
	e := &EtcdCert{
		subject: pkix.Name{
			Country:            []string{config.SubjectConfig.Country},
			Province:           []string{config.SubjectConfig.Province},
			Organization:       []string{config.SubjectConfig.Organization},
			OrganizationalUnit: []string{config.SubjectConfig.OrganizationalUnit},
			CommonName:         config.SubjectConfig.CommonName,
		},
		config: config,
	}
	return e
}

func (e *EtcdCert) generateCACertificate() (*x509.Certificate, *rsa.PrivateKey, error) {

	// 生成私钥
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	ips := make([]net.IP, 0)
	for _, dns := range e.config.DnsList {
		ips = append(ips, net.ParseIP(dns))
	}

	// 构建证书模板
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               e.subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(e.config.Expire),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// 使用模板和私钥生成证书
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	// 保存私钥到文件
	privateKeyFile, err := os.Create(e.config.CaKeyFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer privateKeyFile.Close()
	keyByte := x509.MarshalPKCS1PrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyByte}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return nil, nil, err
	}

	// 保存证书到文件
	certFile, err := os.Create(e.config.CaCertFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer certFile.Close()

	certPEM := &pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	if err := pem.Encode(certFile, certPEM); err != nil {
		return nil, nil, err
	}

	return &template, privateKey, nil
}

func (e *EtcdCert) generateEtcdCertificate(caCert *x509.Certificate, caPrivateKey *rsa.PrivateKey) error {

	ips := make([]net.IP, 0)
	for _, dns := range e.config.DnsList {
		ips = append(ips, net.ParseIP(dns))
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	// 构建etcd证书签署请求 (CSR)
	_, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject:     e.subject,
		IPAddresses: ips,
		DNSNames:    e.config.DnsList,
	}, pk)
	if err != nil {
		return err
	}

	// 使用CA证书和私钥对etcd CSR进行签名，生成etcd证书
	derBytes, err := x509.CreateCertificate(rand.Reader, &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               e.subject,
		IPAddresses:           ips,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}, caCert, &pk.PublicKey, caPrivateKey)
	if err != nil {
		return err
	}

	etcdCertFile, err := os.Create(e.config.EtcdCertFilePath)
	if err != nil {
		return err
	}
	err = pem.Encode(etcdCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return err
	}

	etcdPrivateKeyFile, err := os.Create(e.config.EtcdKeyFilePath)
	if err != nil {
		return err
	}
	err = pem.Encode(etcdPrivateKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	if err != nil {
		return err
	}

	return nil
}

func (e *EtcdCert) GenerateEtcdCert() error {
	caCert, caPrivateKey, err := e.generateCACertificate()
	if err != nil {
		return err
	}
	err = e.generateEtcdCertificate(caCert, caPrivateKey)
	if err != nil {
		return err
	}
	return nil
}
