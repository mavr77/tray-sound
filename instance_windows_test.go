//go:build windows

package main

import (
	"testing"

	"golang.org/x/sys/windows"
)

func TestAlreadyRunning(t *testing.T) {
	first, err := alreadyRunning()
	if err != nil {
		t.Fatal(err)
	}
	if first {
		t.Skip("tray-sound is already running on this machine")
	}
	defer func() {
		_ = windows.CloseHandle(instanceMutex)
		instanceMutex = 0
	}()

	saved := instanceMutex
	second, err := alreadyRunning()
	if err != nil {
		t.Fatal(err)
	}
	if !second {
		t.Fatal("second check should report already running")
	}
	if instanceMutex != saved {
		t.Fatal("second check must not replace the held mutex")
	}
}
