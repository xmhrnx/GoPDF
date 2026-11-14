package studio

import "gopdf/internal/pdf"

// sealLevel orders signature states from best to worst.
type sealLevel int

const (
	sealValid      sealLevel = iota // intact, signer's identity verified
	sealUnverified                  // intact, identity not verified (self-signed)
	sealEmpty                       // a field nobody signed yet
	sealBroken                      // the signed content changed
)

// sealStatus describes one signature: how good it is, a headline and a
// sentence that explains it.
type sealStatus struct {
	level    sealLevel
	headline string
	detail   string
}

// statusOf judges f among all the document's signatures: changes after a
// signature are fine when a later intact signature covers the whole file,
// since they are then other signatures.
func statusOf(f pdf.SignatureField, all []pdf.SignatureField) sealStatus {
	switch {
	case !f.Signed:
		return sealStatus{sealEmpty, "Empty signature field", "This field is waiting to be signed."}
	case !f.Intact:
		return sealStatus{sealBroken, "Signature is invalid", "The document has been changed since " + f.Name + " signed it, or the signature is damaged."}
	}
	later := ""
	if !f.Whole && !wholeCovered(all) {
		later = " The file was changed after signing, in ways the signature doesn't cover."
	}
	if !f.Trusted {
		who := "the signer's certificate is self-signed"
		if !f.SelfSig {
			who = "the certificate's issuer isn't trusted on this computer"
		}
		return sealStatus{sealUnverified, "Signed by " + f.Name + ", identity not verified",
			"Nothing has changed since it was signed, but " + who + "." + later}
	}
	return sealStatus{sealValid, "Signed by " + f.Name + ", valid", "Nothing has changed since it was signed, and the signer's identity is verified." + later}
}

// wholeCovered reports whether an intact signature covers the whole file.
func wholeCovered(all []pdf.SignatureField) bool {
	for _, f := range all {
		if f.Signed && f.Intact && f.Whole {
			return true
		}
	}
	return false
}

// worstStatus is the document's overall signature state.
func worstStatus(all []pdf.SignatureField) (sealLevel, int) {
	worst, signed := sealValid, 0
	for _, f := range all {
		if !f.Signed {
			continue
		}
		signed++
		if l := statusOf(f, all).level; l > worst {
			worst = l
		}
	}
	return worst, signed
}

// sealColor and sealIcon show a level.
func sealColor(l sealLevel) string {
	switch l {
	case sealValid:
		return "theme.success"
	case sealBroken:
		return "theme.danger"
	default:
		return "theme.warning"
	}
}

func sealIcon(l sealLevel) string {
	switch l {
	case sealValid:
		return "shield-check"
	case sealBroken:
		return "shield-x"
	default:
		return "shield-alert"
	}
}
