package recent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddLoad(t *testing.T) {
	dir := t.TempDir()
	s := &Store{file: filepath.Join(dir, "cfg", "recent.json")}
	a, b := filepath.Join(dir, "a.pdf"), filepath.Join(dir, "b.pdf")
	for _, p := range []string{a, b} {
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	list, _ := s.Add(nil, Entry{Path: a})
	list, _ = s.Add(list, Entry{Path: b})
	if _, err := s.Add(list, Entry{Path: a}); err != nil {
		t.Fatal(err)
	}
	os.Remove(b)
	got, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Path != a {
		t.Fatalf("got %+v", got)
	}
}
