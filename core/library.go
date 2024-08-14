package core

func LoadLib(path string) (uintptr, error) {
	return openLibrary(path)
}
