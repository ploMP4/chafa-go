//go:build windows

package chafa

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func loadLibrary() (uintptr, error) {
	libPath, err := extractEmbeddedLibrary()
	if err == nil {
		handle, loadErr := windows.LoadLibrary(libPath)
		if loadErr == nil {
			return uintptr(handle), nil
		}
		fmt.Printf("Warning: Failed to load embedded library: %v\n", loadErr)
	} else {
		fmt.Printf("Warning: Failed to extract embedded library: %v\n", err)
	}

	handle, err := windows.LoadLibrary("libchafa.dll")
	if err != nil {
		return 0, fmt.Errorf("%s library not found: %w", libName, err)
	}
	return uintptr(handle), nil
}
