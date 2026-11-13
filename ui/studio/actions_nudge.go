package studio

import (
	"time"

	"gopdf/internal/geom"
)

// Nudge steps, in points, and how long arrow presses keep adding to the
// same undo step.
const (
	nudgeStep      = 1.0
	nudgeLargeStep = 10.0
	nudgeGroup     = 1500 * time.Millisecond
)

// nudgeState remembers the last arrow-key move, so a run of presses
// undoes in one step.
type nudgeState struct {
	id   int64
	last time.Time
}

// nudge moves the selected markup by (dx, dy) steps; large uses the big
// step (Shift held).
func (s *Studio) nudge(dx, dy float64, large bool) {
	a, ok := s.selectedAnnot()
	if !ok || !s.canEdit() {
		return
	}
	step := nudgeStep
	if large {
		step = nudgeLargeStep
	}
	if s.nudging.id != a.ID || time.Since(s.nudging.last) > nudgeGroup {
		s.doc.BeginChange()
	}
	s.nudging = nudgeState{id: a.ID, last: time.Now()}
	a = a.Clone()
	a.Translate(geom.Point{X: dx * step, Y: dy * step})
	s.keepOnPage(&a)
	s.doc.ReplaceAnnotation(a)
}
