// Package history is a bounded undo/redo stack of immutable states.
package history

// Stack keeps the states before each change (for undo) and the states
// undone (for redo).
type Stack[T any] struct {
	undo, redo []T
	limit      int
}

// New returns a stack that remembers at most limit changes.
func New[T any](limit int) *Stack[T] { return &Stack[T]{limit: limit} }

// Push records before, the state preceding a new change, and forgets
// anything that could be redone.
func (s *Stack[T]) Push(before T) {
	s.undo = append(s.undo, before)
	if len(s.undo) > s.limit {
		s.undo = s.undo[1:]
	}
	s.redo = nil
}

// Undo returns the state to restore, given the current one.
func (s *Stack[T]) Undo(current T) (T, bool) {
	return move(&s.undo, &s.redo, current)
}

// Redo returns the state undone last, given the current one.
func (s *Stack[T]) Redo(current T) (T, bool) {
	return move(&s.redo, &s.undo, current)
}

func move[T any](from, to *[]T, current T) (T, bool) {
	var zero T
	if len(*from) == 0 {
		return zero, false
	}
	last := (*from)[len(*from)-1]
	*from = (*from)[:len(*from)-1]
	*to = append(*to, current)
	return last, true
}

// CanUndo reports whether there is a change to undo.
func (s *Stack[T]) CanUndo() bool { return len(s.undo) > 0 }

// CanRedo reports whether there is a change to redo.
func (s *Stack[T]) CanRedo() bool { return len(s.redo) > 0 }
