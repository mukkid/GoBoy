package GoBoy

import (
	"image"
	"image/png"
	"log"
	"os"
	"testing"
)

// dumpPng will save PNG file to filename of image
func dumpPng(filename string, image *image.RGBA) {
	// Save image for testing
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, image); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func TestVideo(t *testing.T) {
	mem := &GBMem{}

	// write grid tile data
	addr := uint16(VRAMTilePattern)
	mem.write(addr, 0xff)
	mem.write(addr+1, 0xff)
	for i := uint16(2); i < 14; i += 2 {
		mem.write(addr+i, 0xc0)
		mem.write(addr+i+1, 0x03)
	}
	mem.write(addr+14, 0xff)
	mem.write(addr+15, 0xff)

	// write all 0s tile map
	var i uint8 = 0
	for addr := uint16(VRAMBackgroundMap); addr < uint16(VRAMBackgroundMapEnd); addr++ {
		mem.write(addr, 0)
		i = i + 1
	}

	// allocate image
	bgImage := image.NewRGBA(image.Rect(0, 0, 256, 256))

	// drawImage
	bgImage = drawBackground(bgImage, mem)

	// save image for manual inspection
	// FIXME: add actual assertion
	dumpPng("image.png", bgImage)
}
