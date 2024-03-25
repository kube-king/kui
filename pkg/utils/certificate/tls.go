package certificate

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"time"
)

// 初始化证书配置
func initTlsConfig(tlsConfig *TLSConfig) {
	//if tlsConfig.Subject.Country == "" {
	//	tlsConfig.Subject.Country = config.DefaultCountry
	//}
	//if tlsConfig.Subject.Province == "" {
	//	tlsConfig.Subject.Province = config.DefaultProvince
	//}
	//if tlsConfig.Subject.Organization == "" {
	//	tlsConfig.Subject.Organization = config.DefaultOrganization
	//}
	//if tlsConfig.Subject.OrganizationalUnit == "" {
	//	tlsConfig.Subject.OrganizationalUnit = config.DefaultOrganizationUnit
	//}
}

// GenTLSFile 生成tls证书文件
func GenTLSFile(tlsConfig *TLSConfig, expire time.Duration) error {
	tls, err := GenTLS(tlsConfig, expire)
	if err != nil {
		return errors.New(err.Error())
	}

	keyFile, err := os.OpenFile(tlsConfig.KeyOutPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return errors.New(err.Error())
	}
	defer keyFile.Close()
	keyFile.WriteString(tls.Key)

	certFile, err := os.OpenFile(tlsConfig.CertOutPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return errors.New(err.Error())
	}
	defer certFile.Close()
	certFile.WriteString(tls.Cert)

	return nil
}

// GenTLS 生成tls 证书
func GenTLS(tlsConfig *TLSConfig, expire time.Duration) (*TLSContext, error) {

	certBuff := new(bytes.Buffer)
	keyBuff := new(bytes.Buffer)

	tlsContext := &TLSContext{}

	initTlsConfig(tlsConfig)
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)

	// 配置Subject 信息
	subject := pkix.Name{
		Country:            []string{tlsConfig.Subject.Country},
		Province:           []string{tlsConfig.Subject.Province},
		Organization:       []string{tlsConfig.Subject.Organization},
		OrganizationalUnit: []string{tlsConfig.Subject.OrganizationalUnit},
		CommonName:         tlsConfig.Subject.CommonName,
	}

	// certificate 证书模板
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(expire),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		DNSNames: []string{tlsConfig.Subject.CommonName},
	}

	pk, _ := rsa.GenerateKey(rand.Reader, 2048)

	// 创建公钥
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	// certOut, _ := os.Create(tlsConfig.OutPath + "server.pem")
	err := pem.Encode(certBuff, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return nil, err
	}
	tlsContext.Cert = certBuff.String()
	// 创建私钥
	// keyOut, _ := os.Create(tlsConfig.OutPath + "server.key")
	err = pem.Encode(keyBuff, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	if err != nil {
		return nil, err
	}
	tlsContext.Key = keyBuff.String()
	return tlsContext, nil
}
