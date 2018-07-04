package main

import (
	"flag"
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// 256x256 is written to in total but only 160x144 is visible.
const (
	screenWidth   = 256
	screenHeight  = 256
	visibleWidth  = 160
	visibleHeight = 144
)

var keyNames = map[ebiten.Key]string{
	ebiten.KeyDown:    "Down",
	ebiten.KeyLeft:    "Left",
	ebiten.KeyRight:   "Right",
	ebiten.KeyUp:      "Up",
	ebiten.KeySpace:   "Space",
	ebiten.KeyControl: "Ctrl",
}

// global emulation state
var Gb *GameBoy

// getKeys polls for keys defined in keyNames
func getKeys() []string {
	var pressed = []string{}
	for key, name := range keyNames {
		if ebiten.IsKeyPressed(key) {
			pressed = append(pressed, name)
		}
	}
	return pressed
}

// update is the main drawing function
func update(screen *ebiten.Image) error {

	// draw background
	Gb.image = drawBackground(Gb.image, Gb.mainMemory)

	// replace pixels on screen
	screen.ReplacePixels(Gb.image.Pix)

	pressed := getKeys()

	if ebiten.IsRunningSlowly() {
		return nil
	}

	str := fmt.Sprintf("FPS: %f, Keys: %v", ebiten.CurrentFPS(), pressed)
	ebitenutil.DebugPrint(screen, str)

	return nil
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

	// Initialize joypad values
	Gb.mainMemory.ioregs[0] = 0xff

	// step through pc
	go func() {
		for {
			Gb.handleInterrupt()
			Gb.Step()
		}
	}()

	go Gb.incrementLY()

	// setup update loop
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "GoBoy"); err != nil {
		log.Fatal(err)
	}
}
