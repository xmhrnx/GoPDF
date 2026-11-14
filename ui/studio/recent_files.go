package studio

import (
	"time"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/settings"

	"gopdf/internal/recent"
)

// appFolder names the folder digital IDs are kept in, beside the
// settings.
const appFolder = "gopdf"

// recentFiles is the list on the welcome screen.
type recentFiles struct {
	list []recent.Entry
}

func newRecentFiles() *recentFiles { return &recentFiles{} }

// load reads the list on a background goroutine, dropping files that
// are gone.
func (r *recentFiles) load(s *Studio) {
	go func() {
		list := recent.Existing(settings.Value(keyRecent, []recent.Entry(nil)))
		s.win.Post(func() { r.list = list })
	}()
}

// add records path as just opened.
func (r *recentFiles) add(s *Studio, path string, pages int) {
	if path == "" {
		return
	}
	r.list = recent.Add(r.list, recent.Entry{Path: path, Opened: time.Now(), Pages: pages})
	list := r.list
	go func() {
		if err := settings.Save(keyRecent, list); err != nil {
			s.win.Post(func() { s.toasts.Show(giode.ToastWarning, "Couldn't save recent files: "+err.Error()) })
		}
	}()
}
