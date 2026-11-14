package pdf

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"fmt"
	"time"

	pdflib "github.com/digitorus/pdf"
	"github.com/digitorus/pdfsign/sign"

	"gopdf/internal/annot"
	"gopdf/internal/digitalid"
	"gopdf/internal/geom"
)

// Seal signs data with signer, showing mark (a signature annotation on
// page, in display space) as the signature field's appearance. The
// signature is appended as an incremental update, so earlier signatures
// stay valid. password opens an encrypted document; it is ignored
// otherwise.
func Seal(data []byte, page int, info PageInfo, mark annot.Annotation, signer *digitalid.Signer, password string) ([]byte, error) {
	rdr, err := pdflib.NewReaderEncrypted(bytes.NewReader(data), int64(len(data)), func() string { return password })
	if err != nil {
		return nil, fmt.Errorf("read for signing: %w", err)
	}
	user := info.ToUser(mark.Box())
	reason, location := "", ""
	if mark.Seal != nil {
		reason, location = mark.Seal.Reason, mark.Seal.Location
	}
	sd := sign.SignData{
		Signature: sign.SignDataSignature{
			CertType: sign.ApprovalSignature,
			Info: sign.SignDataSignatureInfo{
				Name:        signer.Info.Identity.Name,
				Reason:      reason,
				Location:    location,
				ContactInfo: signer.Info.Identity.Email,
				Date:        time.Now(),
			},
		},
		Signer:            signer.Key,
		DigestAlgorithm:   crypto.SHA256,
		Certificate:       signer.Cert,
		CertificateChains: [][]*x509.Certificate{append([]*x509.Certificate{signer.Cert}, signer.Chain...)},
		SubFilter:         sign.SubFilterETSICAdESDetached,
		Appearance: sign.Appearance{
			Visible:    true,
			Page:       uint32(page + 1),
			LowerLeftX: user.Min.X, LowerLeftY: user.Min.Y,
			UpperRightX: user.Max.X, UpperRightY: user.Max.Y,
			Renderer: func(_ *sign.SignContext, r [4]float64) ([]byte, error) {
				return sealAppearance(mark, info.Rotate, r[2]-r[0], r[3]-r[1]), nil
			},
		},
	}
	var out bytes.Buffer
	if err := sign.Sign(bytes.NewReader(data), &out, rdr, int64(len(data)), sd); err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}
	return out.Bytes(), nil
}

// sealAppearance draws mark as a form XObject w x h in the field's space.
// The page's rotation is undone inside the form, so the signature reads
// upright however the page is turned.
func sealAppearance(mark annot.Annotation, rotate int, w, h float64) []byte {
	m := mark.Clone()
	m.Translate(geom.Point{}.Sub(m.Box().Min))
	m.Seal = nil
	field := PageInfo{Box: geom.Rect{Max: geom.Point{X: w, Y: h}}, Rotate: rotate}
	content := writeAnnotations([]annot.Annotation{m}, field.displayMatrix(), func(float64) string { return "" })
	var b bytes.Buffer
	fmt.Fprintf(&b, "<<\n/Type /XObject\n/Subtype /Form\n/BBox [0 0 %s %s]\n/Matrix [1 0 0 1 0 0]\n/Resources << >>\n/FormType 1\n/Length %d\n>>\nstream\n",
		num(w), num(h), len(content))
	b.Write(content)
	b.WriteString("\nendstream\n")
	return b.Bytes()
}
