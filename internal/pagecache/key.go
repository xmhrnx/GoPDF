// Package pagecache renders pages on a background goroutine and keeps a
// bounded set of the results, so scrolling through a long document
// neither blocks the UI nor holds every page in memory.
package pagecache

// Key identifies one rendering: document version, page and resolution.
type Key struct {
	Version int
	Page    int
	DPI     int
}
