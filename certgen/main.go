package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"math/big"
	unsecure_rand "math/rand"
	"os"
	"time"
)

// GenerateRSA returns a RSA private key with the given key length.
func GenerateRSA(bitSize int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bitSize)
}

// MakeBasicCert returns a basic x509 certificate with minimum fields set to sensible defaults. It is
// expected that users will further modify the certificate before it is used.
func MakeBasicCert() *x509.Certificate {
	//Use a different random number generator so we dont leak any state of crypto/rand
	//Who cares if the serial number is predictable, we know when the cert generated anyway
	//through NotBefore.
	unsecure_rand.Seed(time.Now().Unix())
	return makeBasicCert(time.Now())
}

func makeBasicCert(now time.Time) *x509.Certificate {
	//Make a subjectKeyId. There are no security requirements for this field, but the
	//more statistically distributed it is the better it can be used.
	subjectKeyNum := uint64(unsecure_rand.Int63())
	var subjectKeyBytes = make([]byte, 16)
	binary.PutUvarint(subjectKeyBytes, subjectKeyNum)

	return &x509.Certificate{
		SerialNumber: big.NewInt(int64(unsecure_rand.Int63())),
		Subject: pkix.Name{
			Country:            []string{"U.S"},
			Organization:       []string{"Acme Co."},
			OrganizationalUnit: []string{"Acme Co." + "U"},
		},
		NotBefore:    now,
		NotAfter:     now.AddDate(0, 6, 0), //6 month expiry
		SubjectKeyId: subjectKeyBytes[:5],
	}
}

// MakeBasicServerCert returns a basic x509 certificate with minimum fields necessary to act as a TLS / other server.
// It is expected the caller will set their own fields, like country / subject using SetDetails().
func MakeBasicServerCert() *x509.Certificate {
	cert := MakeBasicCert()
	makeBasicServerCert(cert)
	return cert
}

func makeBasicServerCert(cert *x509.Certificate) {
	cert.IsCA = false
	cert.BasicConstraintsValid = true
	cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	cert.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
}

// MakeBasicCA returns a basic x509 certificate with minimum fields necessary to act as a CA certificate.
// It is expected the caller will set their own fields, like country / subject.
func MakeBasicCA() *x509.Certificate {
	cert := MakeBasicCert()
	makeBasicCA(cert)
	return cert
}

func makeBasicCA(cert *x509.Certificate) {
	cert.IsCA = true
	cert.BasicConstraintsValid = true
	cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	cert.KeyUsage |= x509.KeyUsageCertSign
}

// SetDetails sets the human-readable details of a certificate. Use this after generating a certificate with an above method.
func SetDetails(cert *x509.Certificate, country, organisation, organisationUnit string) {
	cert.Subject = pkix.Name{
		Country:            []string{country},
		Organization:       []string{organisation},
		OrganizationalUnit: []string{organisationUnit},
	}
}

// FullCert represents a valid certificate and private key combination.
type FullCert struct {
	Cert     *x509.Certificate
	DerBytes []byte
	Key      *rsa.PrivateKey
}

// MakeCertKeyPair generates a strong private key and signs the cert with it.
func MakeCertKeyPair(cert *x509.Certificate) (*FullCert, error) {
	var ret FullCert
	priv, err := GenerateRSA(2048)
	if err != nil {
		return nil, err
	}
	ret.Key = priv

	caBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}
	ret.DerBytes = caBytes
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return nil, err
	}
	ret.Cert = caCert
	return &ret, nil
}

// WritePrivateCertToFile writes a full certificate of yours to a set of files, so it can be later loaded.
// passing "" to any of the paths omits that file's generation.
//
func WritePrivateCertToFile(derFile, certPEMFile, keyPEMFile, keyPKCSFile string, cert *FullCert) error {
	if derFile != "" {
		certCerFile, err := os.Create(derFile)
		if err != nil {
			return err
		}
		certCerFile.Write(cert.DerBytes)
		certCerFile.Close()
	}

	if certPEMFile != "" {
		certFile, err := os.Create(certPEMFile)
		if err != nil {
			return err
		}
		pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert.DerBytes})
		certFile.Close()
	}

	if keyPEMFile != "" {
		keyFile, err := os.Create(keyPEMFile)
		if err != nil {
			return err
		}
		pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(cert.Key)})
		keyFile.Close()
	}

	if keyPKCSFile != "" {
		keyFile, err := os.Create(keyPKCSFile)
		if err != nil {
			return err
		}
		privBytes := x509.MarshalPKCS1PrivateKey(cert.Key)
		keyFile.Write(privBytes)
		keyFile.Close()
	}

	return nil
}

// MakeDerivedCertKeyPair generates a strong private key and signs the cert with it.
// It then signs the cert with the CA cert provided, such that the generated cert
// can be proven to be associated with the CA cert.
func MakeDerivedCertKeyPair(ca *FullCert, cert *x509.Certificate) (*FullCert, error) {
	var ret FullCert
	priv, err := GenerateRSA(2048)
	if err != nil {
		return nil, err
	}
	ret.Key = priv

	caBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.Cert, &priv.PublicKey, ca.Key)
	if err != nil {
		return nil, err
	}
	ret.DerBytes = caBytes
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return nil, err
	}
	ret.Cert = caCert
	return &ret, nil
}

func prompt(question string) string {
	fmt.Print(question + ": ")
	var input string
	fmt.Scanln(&input)
	return input
}

func main() {
	org := prompt("What is the name of your organisation/server/group?")
	dns := prompt("What is the domain your server is hosted on?")

	fmt.Print("Generating base (CA) certificate...")
	ca := MakeBasicCA()
	SetDetails(ca, "U.S", org, "Certificate Authority")
	caPair, err := MakeCertKeyPair(ca)
	if err != nil {
		fmt.Println("Error!\n: ", err)
		os.Exit(1)
	}
	err = WritePrivateCertToFile("", "/etc/subshard/ca.pem", "/etc/subshard/ca.key.pem", "", caPair)
	if err != nil {
		fmt.Println("Error!\n: ", err)
		os.Exit(1)
	}
	fmt.Println("DONE.")
	fmt.Println("Written to /etc/subshard/ca.pem & /etc/subshard/ca.key.pem")

	fmt.Print("Generating base (CA) certificate...")
	server := MakeBasicServerCert()
	server.DNSNames = []string{dns}
	server.Subject.CommonName = dns
	SetDetails(server, "U.S", org, "Subshard")
	server.Issuer = ca.Subject
	fullSubjectCert, err := MakeDerivedCertKeyPair(caPair, server)
	if err != nil {
		fmt.Println("Error!\n: ", err)
		os.Exit(1)
	}
	err = WritePrivateCertToFile("", "/etc/subshard/cert.pem", "/etc/subshard/key.pem", "", fullSubjectCert)
	if err != nil {
		fmt.Println("Error!\n: ", err)
		os.Exit(1)
	}
	fmt.Println("DONE.")
	fmt.Println("Written to /etc/subshard/cert.pem & /etc/subshard/key.pem")

	// out, err := exec.Command("sh", "-c", "cat /etc/subshard/cert.server.pem /etc/subshard/ca.pem > /etc/subshard/cert.pem").Output()
	// fmt.Println(string(out))
	// if err != nil {
	// 	fmt.Println("Error!\n: ", err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Written combined certificate to /etc/subshard/cert.pem")
}
