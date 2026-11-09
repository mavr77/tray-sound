//go:build windows

package main

import (
	_ "embed"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Sounds are baked into the exe at build time. Must be PCM WAV.

//go:embed sounds/sound1.wav
var sound1Wav []byte // alert is active

//go:embed sounds/sound2.wav
var sound2Wav []byte // no alert

var (
	winmm         = windows.NewLazySystemDLL("winmm.dll")
	procPlaySound = winmm.NewProc("PlaySoundW")
)

const (
	sndAsync     = 0x0001 // return immediately, play in background
	sndNoDefault = 0x0002 // no fallback beep if the sound fails
	sndMemory    = 0x0004 // first argument points to WAV data in memory
)

// playSound plays WAV bytes asynchronously. A new call stops the previous one.
// The slice must stay alive while playing; package-level vars are safe.
func playSound(wav []byte) {
	if len(wav) == 0 {
		return
	}
	_, _, _ = procPlaySound.Call(
		uintptr(unsafe.Pointer(&wav[0])),
		0,
		sndMemory|sndAsync|sndNoDefault,
	)
}
