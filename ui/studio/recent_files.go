package studio

import (
	"time"

	"github.com/xmhrnx/giode"

	"gopdf/internal/recent"
)

// appFolder names the configuration folder.
const appFolder = "gopdf"

// recentFiles is the list on the welcome screen and in File > Open Recent.
type recentFiles struct {
	store *recent.Store
	list  []recent.Entry
}

func newRecentFiles() *recentFiles {
	store, err := recent.NewStore(appFolder)
	if err != nil {
		return &recentFiles{} // no config folder: the list stays empty
	}
	return &recentFiles{store: store}
}

// load reads the list on a background goroutine.
func (r *recentFiles) load(s *Studio) {
	if r.store == nil {
		return
	}
	go func() {
		list, err := r.store.Load()
		s.win.Post(func() {
			if err != nil {
				s.toasts.Show(giode.ToastWarning, "Couldn't read recent files: "+err.Error())
				return
			}
			r.list = list
		})
	}()
}

// add records path as just opened.
func (r *recentFiles) add(s *Studio, path string, pages int) {
	if r.store == nil || path == "" {
		return
	}
	list, err := r.store.Add(r.list, recent.Entry{Path: path, Opened: time.Now(), Pages: pages})
	r.list = list
	if err != nil {
		s.toasts.Show(giode.ToastWarning, "Couldn't save recent files: "+err.Error())
	}
}
