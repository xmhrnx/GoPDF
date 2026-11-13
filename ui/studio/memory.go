package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/internal/annot"
	"gopdf/internal/prefs"
)

// memory is what the app remembers between runs: the last style used
// for each kind of markup, and the signature to place.
type memory struct {
	store *prefs.Store
	prefs prefs.Prefs
}

func newMemory() *memory {
	store, err := prefs.NewStore(appFolder)
	if err != nil {
		return &memory{} // no config folder: remember for this run only
	}
	return &memory{store: store}
}

// load reads the saved choices on a background goroutine.
func (m *memory) load(s *Studio) {
	if m.store == nil {
		return
	}
	go func() {
		p, err := m.store.Load()
		s.win.Post(func() {
			if err != nil {
				s.toasts.Show(giode.ToastWarning, "Couldn't read your saved settings: "+err.Error())
				return
			}
			m.prefs = p
		})
	}()
}

// style returns the remembered style of kind k.
func (m *memory) style(k annot.Kind) annot.Style { return m.prefs.StyleFor(k) }

// setStyle remembers st for k and saves it.
func (m *memory) setStyle(s *Studio, k annot.Kind, st annot.Style) {
	m.prefs.SetStyle(k, st)
	m.save(s)
}

// setSignature remembers the signature to place.
func (m *memory) setSignature(s *Studio, id string) {
	m.prefs.Signature = id
	m.save(s)
}

func (m *memory) save(s *Studio) {
	if m.store == nil {
		return
	}
	p := m.prefs
	go func() {
		if err := m.store.Save(p); err != nil {
			s.win.Post(func() { s.toasts.Show(giode.ToastWarning, "Couldn't save your settings: "+err.Error()) })
		}
	}()
}
