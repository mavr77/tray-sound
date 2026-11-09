//go:build windows

package main

import (
	"errors"

	"golang.org/x/sys/windows"
)

// instanceMutexName is per logon session (Local\), so each logged-on user
// gets one copy; a second copy in the same session exits.
const instanceMutexName = `Local\tray-sound-single-instance`

// instanceMutex is kept open for the lifetime of the process; Windows
// releases it automatically on exit or crash.
var instanceMutex windows.Handle

// alreadyRunning reports whether another copy of the program holds the
// single-instance mutex. On error it returns false so the app still starts.
func alreadyRunning() (bool, error) {
	name, err := windows.UTF16PtrFromString(instanceMutexName)
	if err != nil {
		return false, err
	}
	h, err := windows.CreateMutex(nil, false, name)
	switch {
	case errors.Is(err, windows.ERROR_ALREADY_EXISTS):
		// Opened the existing mutex; close our handle and report.
		_ = windows.CloseHandle(h)
		return true, nil
	case errors.Is(err, windows.ERROR_ACCESS_DENIED):
		// The mutex exists but we cannot open it (e.g. held by an elevated
		// copy in the same session).
		return true, nil
	case err != nil:
		return false, err
	}
	instanceMutex = h
	return false, nil
}
