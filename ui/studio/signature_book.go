package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/internal/signature"
)

// signatureBook holds the saved signatures.
type signatureBook struct {
	store *signature.Store
	list  []signature.Signature
}

func newSignatureBook() *signatureBook {
	store, err := signature.NewStore(appFolder)
	if err != nil {
		return &signatureBook{}
	}
	return &signatureBook{store: store}
}

// load reads the saved signatures on a background goroutine.
func (b *signatureBook) load(s *Studio) {
	if b.store == nil {
		return
	}
	go func() {
		list, err := b.store.Load()
		s.win.Post(func() {
			if err != nil {
				s.toasts.Show(giode.ToastWarning, "Couldn't read your signatures: "+err.Error())
				return
			}
			b.list = list
		})
	}()
}

// active returns the signature to place: the remembered one, or the
// newest.
func (b *signatureBook) active(s *Studio) (signature.Signature, bool) {
	id := s.memory.prefs.Signature
	for _, sig := range b.list {
		if sig.ID == id {
			return sig, true
		}
	}
	if len(b.list) > 0 {
		return b.list[len(b.list)-1], true
	}
	return signature.Signature{}, false
}

// add saves sig and makes it the one to place.
func (b *signatureBook) add(s *Studio, sig signature.Signature) {
	b.list = append(b.list, sig)
	s.memory.setSignature(s, sig.ID)
	b.save(s)
}

// remove deletes the signature with id. Placed copies stay in documents.
func (b *signatureBook) remove(s *Studio, id string) {
	b.list = signature.Without(b.list, id)
	b.save(s)
}

func (b *signatureBook) save(s *Studio) {
	if b.store == nil {
		return
	}
	list := b.list
	go func() {
		if err := b.store.Save(list); err != nil {
			s.win.Post(func() { s.toasts.Show(giode.ToastError, "Couldn't save signatures: "+err.Error()) })
		}
	}()
}
