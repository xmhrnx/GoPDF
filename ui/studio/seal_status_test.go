package studio

import (
	"testing"

	"gopdf/internal/pdf"
)

func TestStatusOf(t *testing.T) {
	self := pdf.SignatureField{Signed: true, Intact: true, Whole: true, SelfSig: true, Name: "Jane"}
	if l := statusOf(self, nil).level; l != sealUnverified {
		t.Fatalf("self-signed = %v", l)
	}
	trusted := self
	trusted.Trusted = true
	if l := statusOf(trusted, nil).level; l != sealValid {
		t.Fatalf("trusted = %v", l)
	}
	broken := self
	broken.Intact = false
	all := []pdf.SignatureField{trusted, broken}
	if l, n := worstStatus(all); l != sealBroken || n != 2 {
		t.Fatalf("worst = %v of %d", l, n)
	}
}
