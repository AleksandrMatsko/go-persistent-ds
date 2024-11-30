package internal

type VersionMachine struct {
	version uint64
}

func (vm *VersionMachine) GetAndIncrementVersion() uint64 {
	curVersion := vm.version
	vm.version = vm.version + 1
	return curVersion
}

func (vm *VersionMachine) GetVersion() uint64 {
	return vm.version
}
