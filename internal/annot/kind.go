// Package annot is the editable markup layer drawn over PDF pages: ink,
// highlights, shapes and text, kept in page display space until it is
// flattened into the file on save.
package annot

// Kind is the type of an annotation.
type Kind int

// The annotation kinds, in toolbar order.
const (
	Pen Kind = iota
	Marker
	Highlight
	Rect
	Ellipse
	Line
	Arrow
	Text
	Whiteout
	Signature
)

var kindNames = [...]string{
	Pen:       "Pen",
	Marker:    "Marker",
	Highlight: "Highlight",
	Rect:      "Rectangle",
	Ellipse:   "Ellipse",
	Line:      "Line",
	Arrow:     "Arrow",
	Text:      "Text",
	Whiteout:  "Whiteout",
	Signature: "Signature",
}

// Kinds lists every kind, for settings stored per kind.
var Kinds = []Kind{Pen, Marker, Highlight, Rect, Ellipse, Line, Arrow, Text, Whiteout, Signature}

// String returns the kind's display name.
func (k Kind) String() string {
	if int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "Annotation"
}

// Freehand reports whether the kind is drawn as a free path of points.
func (k Kind) Freehand() bool { return k == Pen || k == Marker }

// Boxed reports whether the kind is defined by two opposite corners.
func (k Kind) Boxed() bool {
	return k == Highlight || k == Rect || k == Ellipse || k == Whiteout
}

// Placed reports whether the kind is placed with a click rather than
// drawn with a drag.
func (k Kind) Placed() bool { return k == Text || k == Signature }

// Linear reports whether the kind is a straight segment.
func (k Kind) Linear() bool { return k == Line || k == Arrow }
