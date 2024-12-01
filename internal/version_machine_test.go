package internal

import "testing"

func TestVersionMachine_GetAndIncrementVersion(t *testing.T) {
	vm := VersionMachine{
		version: 0,
	}
	vm.GetAndIncrementVersion()
	if vm.version != 1 {
		t.Error("Version should be 1")
	}
}

func TestVersionMachine_GetVersion(t *testing.T) {
	vm := VersionMachine{
		version: 0,
	}
	vm.GetAndIncrementVersion()
	if vm.GetVersion() != 0 {
		t.Error("Version should be 0")
	}
}
