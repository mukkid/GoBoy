package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// 256x256 is written to in total but only 160x144 is visible.
const (
	screenWidth   = 256
	screenHeight  = 256
	visibleWidth  = 160
	visibleHeight = 144
)

// global emulation state
var Gb *GameBoy

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "GoBoy",
		Bounds: pixel.R(0, 0, 160*3, 144*3),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	d := &pixel.PictureData{
		Pix:    make([]color.RGBA, 144*160),
		Stride: 160,
		Rect:   pixel.R(0, 0, 160, 144),
	}

	for i := 0; i < 144*10; i++ {
		d.Pix[i] = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
		d.Pix[144*160-i-1] = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	}
	for !win.Closed() {
		win.Clear(color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff})
		sprite := pixel.NewSprite(pixel.Picture(d), pixel.R(0, 0, 160, 144))
		mat := pixel.IM.Scaled(pixel.ZV, 3)
		mat = mat.Moved(win.Bounds().Center())
		sprite.Draw(win, mat)

		win.Update()
	}
}

func main() {
	// init gameboy
	Gb = &GameBoy{
		Register:         &Register{},
		mainMemory:       &GBMem{cartridge: &GBROM{}},
		interruptEnabled: true,
		image:            image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}

	// load rom from file
	rom_path := flag.String("rom", "", "rom image to load")
	flag.Parse()
	if *rom_path != "" {
		Gb.mainMemory.cartridge.loadROMFromFile(*rom_path)
	}
	fmt.Println(Gb.rom)
	fmt.Println(Gb.regs)

	pixelgl.Run(run)
	for true {
		Gb.handleInterrupts()
		Gb.Step()
	}
}
