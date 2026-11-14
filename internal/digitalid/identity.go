// Package digitalid creates, imports and stores the digital IDs
// (certificate and private key) documents are signed with. Keys are kept
// only in password-protected PKCS#12 files.
package digitalid

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"time"
)

// Identity is who an ID says its owner is.
type Identity struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
}

// Info describes a stored ID without unlocking it.
type Info struct {
	Fingerprint string    `json:"fingerprint"` // SHA-256 of the certificate, hex
	Identity    Identity  `json:"identity"`
	Issuer      string    `json:"issuer"`
	SelfSigned  bool      `json:"selfSigned"`
	NotAfter    time.Time `json:"notAfter"`
}

// Expired reports whether the certificate is past its end date at t.
func (i Info) Expired(t time.Time) bool { return t.After(i.NotAfter) }

// infoOf describes cert.
func infoOf(cert *x509.Certificate) Info {
	sum := sha256.Sum256(cert.Raw)
	id := Identity{Name: cert.Subject.CommonName}
	if len(cert.EmailAddresses) > 0 {
		id.Email = cert.EmailAddresses[0]
	}
	if len(cert.Subject.Organization) > 0 {
		id.Organization = cert.Subject.Organization[0]
	}
	issuer := cert.Issuer.CommonName
	if issuer == "" && len(cert.Issuer.Organization) > 0 {
		issuer = cert.Issuer.Organization[0]
	}
	return Info{
		Fingerprint: hex.EncodeToString(sum[:]),
		Identity:    id,
		Issuer:      issuer,
		SelfSigned:  selfSigned(cert),
		NotAfter:    cert.NotAfter,
	}
}

// selfSigned reports whether cert is its own issuer and signer.
func selfSigned(cert *x509.Certificate) bool {
	if cert.Subject.String() != cert.Issuer.String() {
		return false
	}
	return cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature) == nil
}
