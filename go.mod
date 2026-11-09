module tray-sound

go 1.25.0

require (
	github.com/gogpu/systray v0.3.0
	golang.org/x/sys v0.47.0
)

require (
	github.com/go-webgpu/goffi v0.6.3 // indirect
	github.com/godbus/dbus/v5 v5.2.2 // indirect
)

// Local copy of gogpu/systray patched so the Windows balloon is silent
// (NIIF_NOSOUND), letting our own sound play alone.
replace github.com/gogpu/systray => ./systray
