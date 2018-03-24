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
	drawBackground(Gb.image, Gb.mainMemory)

	pressed := getKeys()

	if ebiten.IsRunningSlowly() {
		return nil
	}

	str := fmt.Sprintf("FPS: %f, Keys: %v", ebiten.CurrentFPS(), pressed)
	ebitenutil.DebugPrint(screen, str)

	// TODO: Main emulation step, this should be rate limited.
	//if Gb != nil {
	//	Gb.Step()
	//}

	return nil
}

func main() {
	// init gameboy
	Gb = &GameBoy{
		Register:         &Register{},
		rom:              &GBROM{},
		mainMemory:       &GBMem{},
		interruptEnabled: true,
		image:            image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}

	// load rom from file
	rom_path := flag.String("rom", "", "rom image to load")
	if *rom_path != "" {
		Gb.rom.loadROMFromFile(*rom_path)
	}

	// setup update loop
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "GoBoy"); err != nil {
		log.Fatal(err)
	}
}
