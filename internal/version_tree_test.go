package internal

import "testing"

func TestNewVersionTreeCreation(t *testing.T) {
	vt := NewVersionTree[int]()

	if len(vt.tree) != 1 {
		t.Errorf("Expected tree length 1, got: %d", len(vt.tree))
	}
	if vt.tree[0].version != 0 {
		t.Errorf("Expected root version 1, got: %d", vt.tree[0].version)
	}
}

func TestVersionTreeUpdate(t *testing.T) {
	vt := NewVersionTree[int]()

	_, err := vt.Update(0)
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}
	_, err = vt.Update(1)
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}
	_, err = vt.Update(2)
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}
	history, err := vt.GetHistory(3)
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}
	expectedHistory := []uint64{0, 1, 2, 3}
	if !equalSlices(history, expectedHistory) {
		t.Errorf("Expected history: %v, got: %v", expectedHistory, history)
	}

	_, err = vt.Update(1)
	if err != nil {
		return
	}
	_, err = vt.Update(4)
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}
	history, err = vt.GetHistory(5)
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}
	expectedHistory = []uint64{0, 1, 4, 5}
	if !equalSlices(history, expectedHistory) {
		t.Errorf("Expected history: %v, got: %v", expectedHistory, history)
	}
}

func TestVersionTree_UpdateError(t *testing.T) {
	vt := NewVersionTree[int]()

	_, err := vt.Update(2)

	if err == nil {
		t.Error("Expected error, but got none")
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
