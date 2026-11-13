package history

import "testing"

func TestUndoRedo(t *testing.T) {
	s := New[int](2)
	cur := 0
	for _, next := range []int{1, 2, 3} {
		s.Push(cur)
		cur = next
	}
	if cur, _ = s.Undo(cur); cur != 2 {
		t.Fatalf("undo = %d, want 2", cur)
	}
	if cur, _ = s.Undo(cur); cur != 1 {
		t.Fatalf("undo = %d, want 1", cur)
	}
	if _, ok := s.Undo(cur); ok {
		t.Fatal("undo past the limit succeeded")
	}
	if cur, _ = s.Redo(cur); cur != 2 {
		t.Fatalf("redo = %d, want 2", cur)
	}
	s.Push(cur)
	if s.CanRedo() {
		t.Fatal("redo survived a new change")
	}
}
