package studio

import "testing"

func TestFormatPageSize(t *testing.T) {
	cases := map[[2]float64]string{
		{612, 792}: "Letter · 8.5 × 11 in",
		{842, 595}: "A4 · 297 × 210 mm",
		{300, 400}: "106 × 141 mm",
	}
	for in, want := range cases {
		if got := formatPageSize(in[0], in[1]); got != want {
			t.Errorf("formatPageSize(%v) = %q, want %q", in, got, want)
		}
	}
}

func TestNextZoom(t *testing.T) {
	if z := nextZoom(1, 1); z != 1.1 {
		t.Fatalf("zoom in from 100%% = %v", z)
	}
	if z := nextZoom(1.05, -1); z != 1 {
		t.Fatalf("zoom out from 105%% = %v", z)
	}
	if z := nextZoom(4, 1); z != maxZoom {
		t.Fatalf("zoom past max = %v", z)
	}
}
