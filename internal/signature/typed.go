package signature

import (
	"errors"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"

	"gopdf/internal/geom"
)

// Glyph layout constants.
const (
	typedPPEM     = 200 // size the outlines are taken at, in font units per em
	curveSteps    = 10  // line segments per Bézier curve
	maxTypedRunes = 80
)

// ErrNoName is returned for a name with nothing to write.
var ErrNoName = errors.New("type your name first")

// FromText writes name in the TrueType or OpenType font fontData and
// returns it as a signature of filled glyph outlines, so it looks the
// same in every PDF viewer without embedding the font.
func FromText(name string, fontData []byte) (Signature, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Signature{}, ErrNoName
	}
	if r := []rune(name); len(r) > maxTypedRunes {
		name = string(r[:maxTypedRunes])
	}
	f, err := sfnt.Parse(fontData)
	if err != nil {
		return Signature{}, err
	}
	fills, err := layout(f, name)
	if err != nil {
		return Signature{}, err
	}
	sig, ok := normalize(nil, fills)
	if !ok {
		return Signature{}, ErrNoName
	}
	return sig, nil
}

// layout returns the outlines of name's glyphs, set on one line.
func layout(f *sfnt.Font, name string) ([][]geom.Point, error) {
	var (
		buf   sfnt.Buffer
		fills [][]geom.Point
		x     fixed.Int26_6
		prev  sfnt.GlyphIndex
	)
	ppem := fixed.I(typedPPEM)
	for i, r := range name {
		g, err := f.GlyphIndex(&buf, r)
		if err != nil {
			return nil, err
		}
		if g == 0 {
			continue // the font has no glyph for r
		}
		if i > 0 {
			if k, err := f.Kern(&buf, prev, g, ppem, font.HintingNone); err == nil {
				x += k
			}
		}
		segs, err := f.LoadGlyph(&buf, g, ppem, nil)
		if err != nil {
			return nil, err
		}
		fills = append(fills, contours(segs, x)...)
		adv, err := f.GlyphAdvance(&buf, g, ppem, font.HintingNone)
		if err != nil {
			return nil, err
		}
		x += adv
		prev = g
	}
	return fills, nil
}

// contours flattens a glyph's segments, shifted right by dx, into closed
// polygons. sfnt's y axis already points down, like display space.
func contours(segs sfnt.Segments, dx fixed.Int26_6) [][]geom.Point {
	var (
		out []geom.Point
		all [][]geom.Point
	)
	pt := func(p fixed.Point26_6) geom.Point {
		return geom.Point{X: float64(p.X+dx) / 64, Y: float64(p.Y) / 64}
	}
	for _, s := range segs {
		switch s.Op {
		case sfnt.SegmentOpMoveTo:
			if len(out) > 2 {
				all = append(all, out)
			}
			out = []geom.Point{pt(s.Args[0])}
		case sfnt.SegmentOpLineTo:
			out = append(out, pt(s.Args[0]))
		case sfnt.SegmentOpQuadTo:
			out = append(out, quad(out[len(out)-1], pt(s.Args[0]), pt(s.Args[1]))...)
		case sfnt.SegmentOpCubeTo:
			out = append(out, cube(out[len(out)-1], pt(s.Args[0]), pt(s.Args[1]), pt(s.Args[2]))...)
		}
	}
	if len(out) > 2 {
		all = append(all, out)
	}
	return all
}

func quad(p0, p1, p2 geom.Point) []geom.Point {
	pts := make([]geom.Point, 0, curveSteps)
	for i := 1; i <= curveSteps; i++ {
		t := float64(i) / curveSteps
		u := 1 - t
		pts = append(pts, geom.Point{
			X: u*u*p0.X + 2*u*t*p1.X + t*t*p2.X,
			Y: u*u*p0.Y + 2*u*t*p1.Y + t*t*p2.Y,
		})
	}
	return pts
}

func cube(p0, p1, p2, p3 geom.Point) []geom.Point {
	pts := make([]geom.Point, 0, curveSteps)
	for i := 1; i <= curveSteps; i++ {
		t := float64(i) / curveSteps
		u := 1 - t
		pts = append(pts, geom.Point{
			X: u*u*u*p0.X + 3*u*u*t*p1.X + 3*u*t*t*p2.X + t*t*t*p3.X,
			Y: u*u*u*p0.Y + 3*u*u*t*p1.Y + 3*u*t*t*p2.Y + t*t*t*p3.Y,
		})
	}
	return pts
}
