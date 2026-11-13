package pagecache

import (
	"container/list"
	"image"
	"sync"
)

// Renderer draws a page; *pdf.Renderer satisfies it.
type Renderer interface {
	Render(page int, dpi float64) (*image.RGBA, error)
}

// Cache holds up to a fixed number of rendered pages, evicting the least
// recently used, and renders missing ones in the background.
type Cache struct {
	mu       sync.Mutex
	limit    int
	items    map[Key]*list.Element
	order    *list.List // front = most recently used
	pending  map[Key]bool
	queue    chan job
	onReady  func()
	onEvict  func(image.Image)
	renderer Renderer
	version  int
}

type entry struct {
	key Key
	img image.Image
}

type job struct {
	key Key
	r   Renderer
}

// New returns a cache of at most limit images. onReady runs on the
// rendering goroutine after each image is stored (to request a redraw);
// onEvict runs for each image dropped (to release GPU uploads).
func New(limit int, onReady func(), onEvict func(image.Image)) *Cache {
	c := &Cache{
		limit: limit, items: map[Key]*list.Element{}, order: list.New(),
		pending: map[Key]bool{}, queue: make(chan job, 256), onReady: onReady, onEvict: onEvict,
	}
	go c.work()
	return c
}

// SetRenderer switches to a new document version. Queued work for older
// versions is skipped.
func (c *Cache) SetRenderer(r Renderer, version int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.renderer, c.version = r, version
}

// Get returns the image for k, queueing it for rendering when absent. ok
// is false until it is ready; fallback is the closest image of the same
// page at another resolution, to show meanwhile (may be nil).
func (c *Cache) Get(k Key) (img image.Image, fallback image.Image, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, hit := c.items[k]; hit {
		c.order.MoveToFront(el)
		return el.Value.(*entry).img, nil, true
	}
	if !c.pending[k] && c.renderer != nil && k.Version == c.version {
		select {
		case c.queue <- job{key: k, r: c.renderer}:
			c.pending[k] = true
		default: // full: asked again next frame
		}
	}
	return nil, c.closest(k), false
}

func (c *Cache) closest(k Key) image.Image {
	var best *entry
	for el := c.order.Front(); el != nil; el = el.Next() {
		e := el.Value.(*entry)
		if e.key.Page != k.Page || e.key.Version != k.Version {
			continue
		}
		if best == nil || abs(e.key.DPI-k.DPI) < abs(best.key.DPI-k.DPI) {
			best = e
		}
	}
	if best == nil {
		return nil
	}
	return best.img
}

func abs(v int) int { return max(v, -v) }

func (c *Cache) work() {
	for j := range c.queue {
		c.mu.Lock()
		stale := j.key.Version != c.version
		c.mu.Unlock()
		var img image.Image
		if !stale {
			if rgba, err := j.r.Render(j.key.Page, float64(j.key.DPI)); err == nil {
				img = rgba
			}
		}
		c.store(j.key, img)
	}
}

func (c *Cache) store(k Key, img image.Image) {
	var evicted []image.Image
	c.mu.Lock()
	delete(c.pending, k)
	if img != nil {
		c.items[k] = c.order.PushFront(&entry{key: k, img: img})
		for c.order.Len() > c.limit {
			old := c.order.Back().Value.(*entry)
			c.order.Remove(c.order.Back())
			delete(c.items, old.key)
			evicted = append(evicted, old.img)
		}
	}
	c.mu.Unlock()
	if c.onEvict != nil {
		for _, e := range evicted {
			c.onEvict(e)
		}
	}
	if img != nil && c.onReady != nil {
		c.onReady()
	}
}
