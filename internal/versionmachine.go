package internal

type VersionMachine struct {
	version int
}

func (vm *VersionMachine) getAndIncrementVersion() int {
	curVersion := vm.version
	vm.version = vm.version + 1
	return curVersion
}
