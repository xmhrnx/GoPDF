package digitalid

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"math/big"
	"strings"
	"time"
)

// Self-signed ID parameters.
const (
	keyBits  = 2048
	validFor = 5 * 365 * 24 * time.Hour
)

// oidDocumentSigning is the document signing extended key usage
// (RFC 9336), which PDF validators look for.
var oidDocumentSigning = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 36}

// ErrNoName is returned when an identity has no name.
var ErrNoName = errors.New("enter your name")

// unlocked is an ID ready to sign with.
type unlocked struct {
	key   *rsa.PrivateKey
	cert  *x509.Certificate
	chain []*x509.Certificate
}

// newSelfSigned makes a key and a certificate for id, signed by itself.
func newSelfSigned(id Identity, now time.Time) (*unlocked, error) {
	id.Name = strings.TrimSpace(id.Name)
	if id.Name == "" {
		return nil, ErrNoName
	}
	key, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		return nil, err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 127))
	if err != nil {
		return nil, err
	}
	subject := pkix.Name{CommonName: id.Name}
	if o := strings.TrimSpace(id.Organization); o != "" {
		subject.Organization = []string{o}
	}
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               subject,
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(validFor),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageContentCommitment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageEmailProtection},
		UnknownExtKeyUsage:    []asn1.ObjectIdentifier{oidDocumentSigning},
		BasicConstraintsValid: true,
	}
	if e := strings.TrimSpace(id.Email); e != "" {
		tmpl.EmailAddresses = []string{e}
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err
	}
	return &unlocked{key: key, cert: cert}, nil
}
