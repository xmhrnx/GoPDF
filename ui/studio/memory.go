package studio

import (
	"maps"
	"sync"
	"time"

	"github.com/xmhrnx/giode"

	"gopdf/internal/annot"
	"gopdf/internal/prefs"
)

// saveDelay gathers quick changes (a color dragged in the picker) into
// one write.
const saveDelay = 400 * time.Millisecond

// memory is what the app remembers between runs: the last style used
// for each kind of markup, and the signature to place.
type memory struct {
	store *prefs.Store
	prefs prefs.Prefs

	mu      sync.Mutex
	pending *prefs.Prefs // the copy the next write saves
	timer   *time.Timer
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

// save writes the prefs shortly, once changes stop coming.
func (m *memory) save(s *Studio) {
	if m.store == nil {
		return
	}
	p := m.prefs
	p.Styles = maps.Clone(m.prefs.Styles) // the UI goroutine keeps changing the original
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pending = &p
	if m.timer != nil {
		m.timer.Reset(saveDelay)
		return
	}
	m.timer = time.AfterFunc(saveDelay, func() { m.flush(s) })
}

func (m *memory) flush(s *Studio) {
	m.mu.Lock()
	p := m.pending
	m.pending = nil
	m.mu.Unlock()
	if p == nil {
		return
	}
	if err := m.store.Save(*p); err != nil {
		s.win.Post(func() { s.toasts.Show(giode.ToastWarning, "Couldn't save your settings: "+err.Error()) })
	}
}
