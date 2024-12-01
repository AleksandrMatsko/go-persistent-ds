package internal

// VersionMachine is a struct to manage data structure version change.
type VersionMachine struct {
	version uint64
}

// GetAndIncrementVersion returns new version.
func (vm *VersionMachine) GetAndIncrementVersion() uint64 {
	curVersion := vm.version
	vm.version = vm.version + 1
	return curVersion
}

// GetVersion returns current version.
func (vm *VersionMachine) GetVersion() uint64 {
	if vm.version == 0 {
		return 0
	} else {
		return vm.version - 1
	}
}
