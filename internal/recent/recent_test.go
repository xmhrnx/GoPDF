package recent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddExisting(t *testing.T) {
	dir := t.TempDir()
	a, b := filepath.Join(dir, "a.pdf"), filepath.Join(dir, "b.pdf")
	for _, p := range []string{a, b} {
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	list := Add(nil, Entry{Path: a})
	list = Add(list, Entry{Path: b})
	list = Add(list, Entry{Path: a})
	if len(list) != 2 || list[0].Path != a {
		t.Fatalf("list = %+v", list)
	}
	os.Remove(b)
	if got := Existing(list); len(got) != 1 || got[0].Path != a {
		t.Fatalf("existing = %+v", got)
	}
}
