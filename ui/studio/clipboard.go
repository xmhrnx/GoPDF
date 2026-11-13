package studio

import "gopdf/internal/annot"

// pasteOffset is how far each paste lands from the previous one, in
// points, so copies don't hide each other.
const pasteOffset = 12.0

// clipboard holds markup copied with Copy or Cut, for this run.
type clipboard struct {
	item   *annot.Annotation
	pastes int // pastes since the last copy, for the cascade offset
}
