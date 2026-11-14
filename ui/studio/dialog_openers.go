package studio

import "fmt"

// The openers show a dialog when a document is open and nothing runs.

func (v *views) openProperties() {
	if v.s.canEdit() {
		v.properties.show()
	}
}

func (v *views) openWatermark() {
	if v.s.canEdit() {
		v.watermark.show()
	}
}

func (v *views) openPageNumbers() {
	if v.s.canEdit() {
		v.pageNumbers.open = true
	}
}

func (v *views) openSecurity() {
	if v.s.canEdit() {
		v.security.show()
	}
}

func (v *views) openExtract() {
	if v.s.canEdit() {
		v.extract.show(fmt.Sprint(v.s.view.current + 1))
	}
}

func (v *views) extractPage(page int) {
	if v.s.canEdit() {
		v.extract.show(fmt.Sprint(page + 1))
	}
}

func (v *views) openSplit() {
	if v.s.canEdit() {
		v.split.show()
	}
}

func (v *views) openShortcuts() { v.shortcuts.open = true }

func (v *views) openSignatureDetails() {
	if v.s.hasDoc() {
		v.sigDetails.showAll()
	}
}
