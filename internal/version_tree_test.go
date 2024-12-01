package internal

import "testing"

func TestNewVersionTreeCreation(t *testing.T) {
	vt := NewVersionTree(0)
	if vt != nil {
		t.Errorf("NewVersionTree() should return nil")
	}

	vt = NewVersionTree(1)
	if len(vt.tree) != 1 {
		t.Errorf("Expected tree length 1, got %d", len(vt.tree))
	}
	if vt.tree[0].version != 1 {
		t.Errorf("Expected root version 1, got %d", vt.tree[1].version)
	}
}

func TestVersionTreeUpdate(t *testing.T) {
	vt := NewVersionTree(1)

	if !vt.Update(1, 2) {
		t.Error("Expected update to succeed")
	}
	if !vt.Update(2, 3) {
		t.Error("Expected update to succeed")
	}
	if !vt.Update(3, 4) {
		t.Error("Expected update to succeed")
	}

	history := vt.GetHistory(4)
	expectedHistory := []uint64{1, 2, 3, 4}
	if !equalSlices(history, expectedHistory) {
		t.Errorf("Expected history %v, got %v", expectedHistory, history)
	}

	if !vt.Update(1, 5) {
		t.Error("Expected update to succeed")
	}
	if !vt.Update(5, 6) {
		t.Error("Expected update to succeed")
	}
	history = vt.GetHistory(6)
	expectedHistory = []uint64{1, 5, 6}
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
