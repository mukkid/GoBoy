package main

import (
	"fmt"
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

	pressed := getKeys()

	if ebiten.IsRunningSlowly() {
		return nil
	}

	str := fmt.Sprintf("FPS: %f, Keys: %v", ebiten.CurrentFPS(), pressed)
	ebitenutil.DebugPrint(screen, str)

	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "GoBoy"); err != nil {
		log.Fatal(err)
	}
}
