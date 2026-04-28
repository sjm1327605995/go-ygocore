//go:build windows

package core

import "golang.org/x/sys/windows"

func openLibrary(name string) (uintptr, error) {
	handle, err := windows.LoadLibrary(name)
	if err != nil {
		return 0, err
	}
	return uintptr(handle), nil
}

func closeLibrary(handle uintptr) error {
	return windows.FreeLibrary(windows.Handle(handle))
}
