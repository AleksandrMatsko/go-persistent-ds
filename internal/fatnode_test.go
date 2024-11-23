package internal

import (
	"strconv"
	"testing"
)

func TestFatNodeVersion(t *testing.T) {
	fn := NewFatNode("root")

	if fn.root.version != 0 {
		t.Error("Node version is incorrect. Expected 0, but was:" + strconv.Itoa(fn.root.version))
	}
}
