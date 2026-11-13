package pagecache

import (
	"image"
	"sync/atomic"
	"testing"
	"time"
)

type fakeRenderer struct{ calls atomic.Int32 }

func (f *fakeRenderer) Render(page int, dpi float64) (*image.RGBA, error) {
	f.calls.Add(1)
	return image.NewRGBA(image.Rect(0, 0, int(dpi), page+1)), nil
}

func waitFor(t *testing.T, c *Cache, k Key) image.Image {
	t.Helper()
	for range 200 {
		if img, _, ok := c.Get(k); ok {
			return img
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("%v never rendered", k)
	return nil
}

func TestCacheEvictsAndFallsBack(t *testing.T) {
	r := &fakeRenderer{}
	var evicted atomic.Int32
	c := New(2, nil, func(image.Image) { evicted.Add(1) })
	c.SetRenderer(r, 1)
	waitFor(t, c, Key{1, 0, 72})
	waitFor(t, c, Key{1, 1, 72})
	waitFor(t, c, Key{1, 2, 72}) // evicts page 0
	if _, _, ok := c.Get(Key{1, 0, 72}); ok || evicted.Load() != 1 {
		t.Fatal("page 0 not evicted")
	}
	_, fb, ok := c.Get(Key{1, 2, 144})
	if ok || fb == nil || fb.Bounds().Dx() != 72 {
		t.Fatalf("fallback = %v, ok %v", fb, ok)
	}
	c.SetRenderer(r, 2)
	if _, fb, _ := c.Get(Key{2, 2, 72}); fb != nil {
		t.Fatal("fallback from an older version")
	}
}
