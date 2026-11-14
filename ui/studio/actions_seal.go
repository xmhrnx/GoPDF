package studio

import (
	"gopdf/internal/annot"
	"gopdf/internal/digitalid"
)

// setSealed turns the selected signature into a digital signature, or
// back into a plain mark. Without a digital ID it opens the ID manager
// first.
func (s *Studio) setSealed(on bool) {
	if !on {
		s.updateSelected(func(a *annot.Annotation) bool {
			a.Seal = nil
			return a.Kind == annot.Signature
		})
		return
	}
	id := s.defaultDigitalID()
	if id == "" {
		s.ui.digitalIDs.show(func(info digitalid.Info) { s.sealWith(info.Fingerprint) })
		return
	}
	s.sealWith(id)
}

// sealWith makes the selected signature a digital one signed with id.
func (s *Studio) sealWith(id string) {
	s.memory.prefs.DigitalID = id
	s.memory.save(s)
	s.updateSelected(func(a *annot.Annotation) bool {
		if a.Kind != annot.Signature {
			return false
		}
		seal := annot.Seal{ID: id}
		if a.Seal != nil {
			seal.Reason, seal.Location = a.Seal.Reason, a.Seal.Location
		}
		a.Seal = &seal
		return true
	})
}

// defaultDigitalID is the ID used last, or the first one stored.
func (s *Studio) defaultDigitalID() string {
	if _, ok := s.ids.info(s.memory.prefs.DigitalID); ok {
		return s.memory.prefs.DigitalID
	}
	if len(s.ids.list) > 0 {
		return s.ids.list[0].Fingerprint
	}
	return ""
}

func (s *Studio) setSealReason(text string) {
	s.updateSelectedLive("reason", func(a *annot.Annotation) bool {
		if a.Seal == nil {
			return false
		}
		a.Seal.Reason = text
		return true
	})
}

func (s *Studio) setSealLocation(text string) {
	s.updateSelectedLive("location", func(a *annot.Annotation) bool {
		if a.Seal == nil {
			return false
		}
		a.Seal.Location = text
		return true
	})
}
