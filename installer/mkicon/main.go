// mkicon writes installer/bell.ico: the app's bell drawing at several sizes,
// stored as PNG frames (supported by Windows Vista and later).
//
//	go run ./installer/mkicon installer/bell.ico
package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

var bellColor = color.RGBA{R: 220, G: 30, B: 30, A: 255}

func main() {
	out := "installer/bell.ico"
	if len(os.Args) > 1 {
		out = os.Args[1]
	}

	sizes := []int{16, 24, 32, 48, 64, 128, 256}
	frames := make([][]byte, len(sizes))
	for i, s := range sizes {
		frames[i] = bell(s)
	}

	var buf bytes.Buffer
	w := func(v any) { _ = binary.Write(&buf, binary.LittleEndian, v) }

	// ICONDIR
	w(uint16(0))
	w(uint16(1)) // type: icon
	w(uint16(len(sizes)))

	// ICONDIRENTRY per frame
	offset := 6 + 16*len(sizes)
	for i, s := range sizes {
		dim := uint8(s)
		if s >= 256 {
			dim = 0 // 0 means 256
		}
		w(dim)       // width
		w(dim)       // height
		w(uint8(0))  // palette colors
		w(uint8(0))  // reserved
		w(uint16(1)) // color planes
		w(uint16(32))
		w(uint32(len(frames[i])))
		w(uint32(offset))
		offset += len(frames[i])
	}
	for _, f := range frames {
		buf.Write(f)
	}

	if err := os.WriteFile(out, buf.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}

// bell draws the same bell shape as generateBellIcon in main.go,
// with a thicker clapper so it stays visible at large sizes.
func bell(size int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	mid := size / 2
	topY, botY := size/4, size*3/4
	clapH := max(2, size/11)
	clapW := max(1, size/22)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if y >= topY && y <= botY {
				p := float64(y-topY) / float64(botY-topY)
				hw := int(float64(size/6) + p*float64(size/4))
				if x >= mid-hw && x <= mid+hw {
					img.SetRGBA(x, y, bellColor)
				}
			}
			if y > botY && y <= botY+clapH && x >= mid-clapW && x <= mid+clapW {
				img.SetRGBA(x, y, bellColor)
			}
		}
	}

	var b bytes.Buffer
	_ = png.Encode(&b, img)
	return b.Bytes()
}
