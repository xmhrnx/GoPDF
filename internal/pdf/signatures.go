package pdf

import (
	"bytes"
	"crypto/x509"
	"time"

	pdflib "github.com/digitorus/pdf"
	"github.com/digitorus/pdfsign/verify"

	"gopdf/internal/geom"
)

// SignatureField is a signature field of a document and what checking
// its signature found.
type SignatureField struct {
	Page    int       // zero-based; -1 for a signature without a place on a page
	Rect    geom.Rect // display space of Page; empty when invisible
	Signed  bool      // false for an empty field waiting to be signed
	Name    string    // the signer, from the certificate
	Email   string
	Org     string
	Issuer  string // who vouches for the signer's certificate
	SelfSig bool   // the signer vouches for themselves
	Reason  string
	Place   string
	Contact string
	When    time.Time // signing time the signer stated; zero when absent
	Intact  bool      // the signed content hasn't changed since
	Trusted bool      // the certificate chains to a root this computer trusts
	Whole   bool      // nothing at all was added to the file after this signature
	Errors  []string  // why Intact or Trusted is false, as the checker put it
}

// ReadSignatures lists the signature fields of data, checking each
// signature offline. pages gives the geometry to place fields with;
// password opens an encrypted document.
func ReadSignatures(data []byte, pages []PageInfo, password string) ([]SignatureField, error) {
	rdr, err := pdflib.NewReaderEncrypted(bytes.NewReader(data), int64(len(data)), func() string { return password })
	if err != nil {
		return nil, err
	}
	opts := verify.DefaultVerifyOptions()
	opts.Password = password
	var out []SignatureField
	for i := 1; i <= rdr.NumPage() && i <= len(pages); i++ {
		annots := rdr.Page(i).V.Key("Annots")
		for j := 0; j < annots.Len(); j++ {
			widget := annots.Index(j)
			field := widget
			if field.Key("FT").IsNull() {
				field = widget.Key("Parent")
			}
			if widget.Key("Subtype").Name() != "Widget" || field.Key("FT").Name() != "Sig" {
				continue
			}
			sf := SignatureField{Page: i - 1, Rect: widgetRect(widget, pages[i-1])}
			if v := field.Key("V"); !v.IsNull() {
				check(&sf, v, data, opts)
			}
			out = append(out, sf)
		}
	}
	return out, nil
}

func widgetRect(widget pdflib.Value, page PageInfo) geom.Rect {
	r := widget.Key("Rect")
	if r.Len() != 4 {
		return geom.Rect{}
	}
	user := geom.RectFromPoints(
		geom.Point{X: r.Index(0).Float64(), Y: r.Index(1).Float64()},
		geom.Point{X: r.Index(2).Float64(), Y: r.Index(3).Float64()},
	)
	if user.Dx() == 0 || user.Dy() == 0 {
		return geom.Rect{}
	}
	return page.ToDisplay(user)
}

// check verifies the signature dictionary v into sf.
func check(sf *SignatureField, v pdflib.Value, data []byte, opts *verify.VerifyOptions) {
	sf.Signed = true
	br := v.Key("ByteRange")
	sf.Whole = br.Len() == 4 && br.Index(2).Int64()+br.Index(3).Int64() == int64(len(data))
	signer, err := verify.VerifySignature(v, bytes.NewReader(data), int64(len(data)), opts)
	if err != nil {
		sf.Errors = append(sf.Errors, err.Error())
	}
	if signer == nil {
		return
	}
	sf.Reason, sf.Place, sf.Contact = signer.Reason, signer.Location, signer.ContactInfo
	sf.Name = signer.Name
	if signer.SignatureTime != nil {
		sf.When = *signer.SignatureTime
	}
	sf.Intact, sf.Trusted = signer.ValidSignature, signer.TrustedIssuer && !signer.RevokedCertificate
	for _, e := range signer.ValidationErrors {
		sf.Errors = append(sf.Errors, e.Error())
	}
	if len(signer.Certificates) > 0 && signer.Certificates[0].Certificate != nil {
		describeSigner(sf, signer.Certificates[0].Certificate)
	}
}

func describeSigner(sf *SignatureField, cert *x509.Certificate) {
	if cn := cert.Subject.CommonName; cn != "" {
		sf.Name = cn
	}
	if len(cert.EmailAddresses) > 0 {
		sf.Email = cert.EmailAddresses[0]
	}
	if len(cert.Subject.Organization) > 0 {
		sf.Org = cert.Subject.Organization[0]
	}
	sf.Issuer = cert.Issuer.CommonName
	sf.SelfSig = cert.Subject.String() == cert.Issuer.String()
}
