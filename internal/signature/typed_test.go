package signature

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFromTextEveryFont(t *testing.T) {
	files, _ := filepath.Glob("../../assets/signature-fonts/*.ttf")
	if len(files) == 0 {
		t.Fatal("no signature fonts")
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		sig, err := FromText("Jane Q. Doe", data)
		if err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		if len(sig.Fills) < 5 || sig.Aspect < 1.5 {
			t.Fatalf("%s: %d outlines, aspect %v", file, len(sig.Fills), sig.Aspect)
		}
		for _, poly := range sig.Fills {
			for _, p := range poly {
				if p.X < 0 || p.X > 1 || p.Y < 0 || p.Y > 1 {
					t.Fatalf("%s: point %v outside the unit box", file, p)
				}
			}
		}
	}
}

func TestFromTextEmpty(t *testing.T) {
	data, _ := os.ReadFile("../../assets/signature-fonts/Sacramento-Regular.ttf")
	if _, err := FromText("   ", data); !errors.Is(err, ErrNoName) {
		t.Fatalf("err = %v", err)
	}
}
