package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/internal/digitalid"
	"gopdf/internal/editor"
)

// digitalIDs are the user's digital IDs and the ones unlocked this run.
type digitalIDs struct {
	store    *digitalid.Store
	list     []digitalid.Info
	unlocked editor.Signers
}

func newDigitalIDs() *digitalIDs {
	ids := &digitalIDs{unlocked: editor.Signers{}}
	if store, err := digitalid.NewStore(appFolder); err == nil {
		ids.store = store
	}
	return ids
}

// load reads the list on a background goroutine.
func (ids *digitalIDs) load(s *Studio) {
	if ids.store == nil {
		return
	}
	go func() {
		list, err := ids.store.List()
		s.win.Post(func() {
			if err != nil {
				s.toasts.Show(giode.ToastWarning, "Couldn't read your digital IDs: "+err.Error())
				return
			}
			ids.list = list
		})
	}()
}

// info returns the stored ID with fingerprint.
func (ids *digitalIDs) info(fingerprint string) (digitalid.Info, bool) {
	for _, i := range ids.list {
		if i.Fingerprint == fingerprint {
			return i, true
		}
	}
	return digitalid.Info{}, false
}

// added records a newly created or imported ID.
func (ids *digitalIDs) added(info digitalid.Info) {
	for i, old := range ids.list {
		if old.Fingerprint == info.Fingerprint {
			ids.list[i] = info
			return
		}
	}
	ids.list = append(ids.list, info)
}

// remove deletes the ID with fingerprint.
func (ids *digitalIDs) remove(s *Studio, fingerprint string) {
	delete(ids.unlocked, fingerprint)
	s.runBackground("Deleting digital ID", func() error { return ids.store.Delete(fingerprint) }, func() {
		for i, info := range ids.list {
			if info.Fingerprint == fingerprint {
				ids.list = append(ids.list[:i], ids.list[i+1:]...)
				break
			}
		}
	})
}
