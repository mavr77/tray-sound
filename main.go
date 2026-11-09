// Tray app that checks the Kyiv air-alert status every minute.
//
//   - active == true:  icon color = "level", plays sound1.wav
//   - active == false: icon green, plays sound2.wav
//   - server unreachable: icon gray, no sound
//
// Build without a console window on Windows:
//
//	go build -ldflags "-H=windowsgui -s -w" -o tray-sound.exe .
package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gogpu/systray"
)

func main() {
	setupLog()

	running, err := alreadyRunning()
	if err != nil {
		log.Println("single-instance check failed:", err)
	}
	if running {
		log.Println("another instance is already running, exiting")
		return
	}

	tray := systray.New()
	menu := systray.NewMenu()

	checkNow := make(chan struct{}, 1)

	menu.Add("Перевірити зараз", func() {
		select {
		case checkNow <- struct{}{}:
		default: // a check is already queued
		}
	})
	menu.Add("Тест: звук тривоги", func() { playSound(sound1Wav) })
	menu.Add("Тест: звук відбою", func() { playSound(sound2Wav) })
	menu.AddSeparator()
	menu.Add("Вихід", func() {
		log.Println("exiting")
		tray.Remove()
		os.Exit(0)
	})

	tray.SetIcon(generateBellIcon(22, colGray)).
		SetTooltip("Перевірка тривоги...").
		SetMenu(menu).
		Show()

	go pollLoop(tray, checkNow)

	log.Println("running")
	if err := tray.Run(); err != nil {
		log.Println("run error:", err)
	}
}

// pollLoop checks right away, then every pollInterval or on "check now".
func pollLoop(tray *systray.SystemTray, checkNow <-chan struct{}) {
	icons := map[color.RGBA][]byte{} // cache generated icons by color
	iconFor := func(c color.RGBA) []byte {
		if b, ok := icons[c]; ok {
			return b
		}
		b := generateBellIcon(22, c)
		icons[c] = b
		return b
	}

	lastKey := "" // empty until the first successful check

	check := func() {
		now := time.Now()
		st, err := fetchAlert()
		if err != nil {
			log.Println("fetch error:", err)
			tray.SetIcon(iconFor(colGray)).
				SetTooltip("Помилка перевірки " + now.Format("15:04:05"))
			return
		}

		key := st.key()
		changed := key != lastKey
		log.Printf("active=%v level=%q changed=%v", st.Active, st.Level, changed)

		tray.SetIcon(iconFor(st.color())).SetTooltip(st.tooltip(now))

		if playSoundEveryCheck || changed {
			if st.Active {
				playSound(sound1Wav)
			} else {
				playSound(sound2Wav)
			}
		}

		// Pop-up only when the state changes (not on the very first check
		// if there is no alert).
		if changed && (lastKey != "" || st.Active) {
			if st.Active {
				tray.ShowNotification("Повітряна тривога",
					st.name()+": Терміново пройдіть до укриття! рівень небезпеки "+map[string]string{"yellow": "жовтий", "red": "червоний"}[st.Level])
			} else {
				tray.ShowNotification("Відбій тривоги", st.name()+": тривоги немає")
			}
		}

		lastKey = key
	}

	check()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			check()
		case <-checkNow:
			check()
		}
	}
}

// setupLog writes logs to %LOCALAPPDATA%\tray-sound\tray-sound.log, because a
// windowsgui build has no console and Program Files is not writable for
// ordinary users. Falls back to the exe's folder if LOCALAPPDATA is unset.
func setupLog() {
	dir := "."
	if exe, err := os.Executable(); err == nil {
		dir = filepath.Dir(exe)
	}
	if base := os.Getenv("LOCALAPPDATA"); base != "" {
		d := filepath.Join(base, "tray-sound")
		if err := os.MkdirAll(d, 0o755); err == nil {
			dir = d
		}
	}
	f, err := os.OpenFile(filepath.Join(dir, "tray-sound.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	log.SetOutput(f)
}

// generateBellIcon draws a simple bell icon of the given color as PNG bytes.
func generateBellIcon(size int, c color.RGBA) []byte {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	mid := size / 2
	topY, botY := size/4, size*3/4

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if y >= topY && y <= botY {
				p := float64(y-topY) / float64(botY-topY)
				hw := int(float64(size/6) + p*float64(size/4))
				if x >= mid-hw && x <= mid+hw {
					img.SetRGBA(x, y, c)
				}
			}
			if y > botY && y <= botY+2 && x >= mid-1 && x <= mid+1 {
				img.SetRGBA(x, y, c)
			}
		}
	}

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}
