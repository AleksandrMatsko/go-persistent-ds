package internal

import "testing"

func TestNewVersionTreeCreation(t *testing.T) {
	vt := NewVersionTree()

	if len(vt.tree) != 1 {
		t.Errorf("Expected tree length 1, got %d", len(vt.tree))
	}
	if vt.tree[0].version != 0 {
		t.Errorf("Expected root version 1, got %d", vt.tree[0].version)
	}
}

func TestVersionTreeUpdate(t *testing.T) {
	vt := NewVersionTree()

	vt.Update(0)
	vt.Update(1)
	vt.Update(2)
	history := vt.GetHistory(3)
	expectedHistory := []uint64{0, 1, 2, 3}
	if !equalSlices(history, expectedHistory) {
		t.Errorf("Expected history %v, got %v", expectedHistory, history)
	}

	vt.Update(1)
	vt.Update(4)
	history = vt.GetHistory(5)
	expectedHistory = []uint64{0, 1, 4, 5}
	if !equalSlices(history, expectedHistory) {
		t.Errorf("Expected history %v, got %v", expectedHistory, history)
	}
}

func equalSlices(a, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
