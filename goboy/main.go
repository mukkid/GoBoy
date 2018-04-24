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
var (
	GbRom     *GBROM
	GbRam     *GBMem
	GbBGImage *image.RGBA
)

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
	drawBackground(GbBGImage, GbRam)

	pressed := getKeys()

	if ebiten.IsRunningSlowly() {
		return nil
	}

	str := fmt.Sprintf("FPS: %f, Keys: %v", ebiten.CurrentFPS(), pressed)
	ebitenutil.DebugPrint(screen, str)

	return nil
}

func main() {

	// init global vars
	GbRom = &GBROM{}
	GbRam = &GBMem{}

	// load rom from file
	rom_path := flag.String("rom", "", "rom image to load")
	if *rom_path != "" {
		GbRom.loadROMFromFile(*rom_path)
	}

	// allocate image buffer
	GbBGImage = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "GoBoy"); err != nil {
		log.Fatal(err)
	}
}
