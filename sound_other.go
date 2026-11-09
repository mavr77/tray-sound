//go:build !windows

package main

// Custom sound is implemented only for Windows.
var (
	sound1Wav []byte
	sound2Wav []byte
)

func playSound(wav []byte) {}
