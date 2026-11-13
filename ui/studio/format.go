package studio

import (
	"fmt"
	"math"
)

// paperSizes are well-known page sizes in points, portrait.
var paperSizes = []struct {
	name string
	w, h float64
}{
	{"Letter", 612, 792},
	{"Legal", 612, 1008},
	{"A3", 842, 1191},
	{"A4", 595, 842},
	{"A5", 420, 595},
}

// formatPageSize describes a page size in points: "A4 · 210 × 297 mm".
func formatPageSize(w, h float64) string {
	mm := fmt.Sprintf("%.0f × %.0f mm", w/72*25.4, h/72*25.4)
	short, long := min(w, h), max(w, h)
	for _, p := range paperSizes {
		if math.Abs(short-p.w) < 3 && math.Abs(long-p.h) < 3 {
			if p.name == "Letter" || p.name == "Legal" {
				return fmt.Sprintf("%s · %g × %g in", p.name, round1(w/72), round1(h/72))
			}
			return p.name + " · " + mm
		}
	}
	return mm
}

func round1(v float64) float64 { return math.Round(v*10) / 10 }

// formatBytes formats a size as "1.4 MB".
func formatBytes(n int) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.0f KB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%d B", n)
	}
}
