//go:build !windows

package main

// Single-instance check is implemented only for Windows.
func alreadyRunning() (bool, error) { return false, nil }
