package internal

type VersionMachine struct {
	version int
}

func (vm *VersionMachine) GetAndIncrementVersion() int {
	curVersion := vm.version
	vm.version = vm.version + 1
	return curVersion
}

func (vm *VersionMachine) GetVersion() int {
	return vm.version
}
