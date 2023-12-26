/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @since: 2023/6/11
  @desc: //TODO
**/

package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

// LoadCertificateTls 加载证书
func LoadCertificateTls(crtPath string, keyPath string) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		log.Println("load certificate failed", err)
		return nil, err
	}
	return &cert, nil
}

// LoadCertificateX509Data 加载证书私钥数据
func LoadCertificateX509Data(crtPath string, keyPath string) (*x509.Certificate, string, error) {
	rootCrtData, err := os.ReadFile(crtPath)
	if err != nil {
		log.Println("load certificate failed", err)
		panic(err)
	}
	rootBlock, _ := pem.Decode(rootCrtData)
	if rootBlock == nil || rootBlock.Type != "CERTIFICATE" {
		panic(err)
	}
	rootCrt, err := x509.ParseCertificate(rootBlock.Bytes)
	if err != nil {
		return nil, "", err
	}
	rootKeyData, err := os.ReadFile(keyPath)
	if err != nil {
		log.Println("load certificate failed", err)
		return nil, "", err
	}
	return rootCrt, string(rootKeyData), err
}

// GenerateCaCertificate 生成Ca证书
func GenerateCaCertificate(rootCrt *tls.Certificate, host string) (*tls.Certificate, error) {
	host = strings.TrimSuffix(host, ":443")
	interCsr := &x509.Certificate{
		Version:      3,
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Province:           []string{"CiProxy"},
			Locality:           []string{"GuangZhou"},
			Organization:       []string{"CiProxy"},
			OrganizationalUnit: []string{"CiProxyHttpsSniff"},
			CommonName:         host,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment | x509.KeyUsageKeyEncipherment,
		EmailAddresses:        []string{"qingqianludao@gmail.com"},
	}
	interKey := generateEccPrivateKey()
	interDer, err := x509.CreateCertificate(rand.Reader, interCsr, interCsr, interKey.Public(), rootCrt.PrivateKey)
	if err != nil {
		panic(err)
	}
	interCertX509, err := x509.ParseCertificate(interDer)
	if err != nil {
		panic(err)
	}
	caCertificate := tls.Certificate{
		Certificate: [][]byte{interDer},
		PrivateKey:  rootCrt.PrivateKey,
		Leaf:        interCertX509,
	}

	return &caCertificate, err
}

// generateEccPrivateKey 生成 ECC 私钥
func generateEccPrivateKey() (key *ecdsa.PrivateKey) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return key
}

// generateRsaPrivateKey 生成 Rsa私钥
func generateRsaPrivateKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return key
}

// saveCert 保存证书
func saveCert(cert *x509.Certificate, fileName string) {
	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	pemData := pem.EncodeToMemory(certBlock)
	if err := os.WriteFile(fileName, pemData, 0644); err != nil {
		log.Println("cert save failed")
		panic(err)
	}
}

func GenerateTlsConfig(host string) (*tls.Config, error) {
	cert, err := LoadCertificateTls("./cert/root.crt", "./cert/private.pem")
	if err != nil {
		//errLog("load root certificate failed", err)
		panic(err)
		return nil, err
	}

	caCertificate, err := GenerateCaCertificate(cert, host)
	//if err != nil {
	//	errLog("load ca certificate failed", err)
	//	panic(err)
	//	return
	//}
	conf := &tls.Config{
		Certificates: []tls.Certificate{*caCertificate},
		//InsecureSkipVerify: true,
		MaxVersion: tls.VersionTLS12,
	}
	return conf, err
}

// savePrivateKey 保存私钥
func savePrivateKey(key *ecdsa.PrivateKey, fileName string) {
	keyDer, err := x509.MarshalECPrivateKey(key)

	keyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyDer,
	}

	keyData := pem.EncodeToMemory(keyBlock)

	if err = os.WriteFile(fileName, keyData, 0777); err != nil {
		log.Println("private key save failed")
		panic(err)
	}

}

// GenerateCert 生成根证书，中级证书， 终端证书
func GenerateCert(fileType string, organization string, country string, province string, locality string, organizationalUnit string, commonName string, dnsDomain string) {

	filePath := "./cert/" + organization + "/"
	println(filePath)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(filePath, 0755)
		if err != nil {
			panic(err)
		}
	}
	// 生成根证书
	rootCsr := &x509.Certificate{
		Version:      3,
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Country:            []string{country},
			Province:           []string{province},
			Locality:           []string{locality},
			Organization:       []string{organization},
			OrganizationalUnit: []string{organizationalUnit},
			CommonName:         commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		MaxPathLenZero:        false,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	rootKey := generateEccPrivateKey()
	rootDer, err := x509.CreateCertificate(rand.Reader, rootCsr, rootCsr, rootKey.Public(), rootKey)
	if err != nil {
		panic(err)
	}
	rootCert, err := x509.ParseCertificate(rootDer)
	if err != nil {
		panic(err)
	}
	saveCert(rootCert, filePath+"root.crt")
	savePrivateKey(rootKey, filePath+"root.key")
	//// 根证书签证中级证书
	//interCsr := &x509.Certificate{
	//	Version:      3,
	//	SerialNumber: big.NewInt(time.Now().Unix()),
	//	Subject: pkix.Name{
	//		Country:            []string{country},
	//		Province:           []string{province},
	//		Locality:           []string{locality},
	//		Organization:       []string{organization},
	//		OrganizationalUnit: []string{organizationalUnit},
	//		CommonName:         commonName,
	//	},
	//	NotBefore:             time.Now(),
	//	NotAfter:              time.Now().AddDate(1, 0, 0),
	//	BasicConstraintsValid: true,
	//	IsCA:                  true,
	//	MaxPathLen:            0,
	//	MaxPathLenZero:        true,
	//	KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	//}
	//interKey := generatePrivateKey()
	//interDer, err := x509.CreateCertificate(rand.Reader, interCsr, rootCert, interKey.Public(), rootKey)
	//if err != nil {
	//	panic(err)
	//}
	//interCert, err := x509.ParseCertificate(interDer)
	//if err != nil {
	//	panic(err)
	//}
	//saveCert(interCert, filePath+"inter.crt")
	//savePrivateKey(interKey, filePath+"inter.key")
	//// 中级证书签证终端证书
	//csr := &x509.Certificate{
	//	Version:      3,
	//	SerialNumber: big.NewInt(time.Now().Unix()),
	//	Subject: pkix.Name{
	//		Country:            []string{country},
	//		Province:           []string{province},
	//		Locality:           []string{locality},
	//		Organization:       []string{organization},
	//		OrganizationalUnit: []string{organizationalUnit},
	//		CommonName:         commonName,
	//	},
	//	DNSNames:              []string{dnsDomain},
	//	NotBefore:             time.Now(),
	//	NotAfter:              time.Now().AddDate(1, 0, 0),
	//	BasicConstraintsValid: true,
	//	IsCA:                  false,
	//	KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	//	ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	//}
	//key := generatePrivateKey()
	//der, err := x509.CreateCertificate(rand.Reader, csr, interCert, key.Public(), interKey)
	//if err != nil {
	//	panic(err)
	//}
	//cert, err := x509.ParseCertificate(der)
	//if err != nil {
	//	panic(err)
	//}
	//saveCert(cert, filePath+"server.crt")
	//savePrivateKey(key, filePath+"server.key")
}
