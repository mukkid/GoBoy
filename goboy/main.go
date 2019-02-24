package main

import (
	"flag"
	"fmt"
	"image"
    "os/signal"
    "syscall"
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

func main() {
	// init gameboy
	Gb = &GameBoy{
		Register:         &Register{},
		mainMemory:       &GBMem{cartridge: &GBROM{}},
		interruptEnabled: true,
		image:            image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}

    d := Debugger {
        gb: Gb,
        breakpoints: make(map[uint16]struct{}),
        paused: true,
    }

	// load rom from file
	rom_path := flag.String("rom", "", "rom image to load")
	flag.Parse()
	if *rom_path != "" {
		Gb.mainMemory.cartridge.loadROMFromFile(*rom_path)
        fmt.Printf("Loaded %s\n", *rom_path)
	}

	// Initialize joypad values
	Gb.mainMemory.ioregs[0] = 0xff

    /* Initialize SIGINT handler */
    debuggersSignal = append(debuggersSignal, &d)
    go sigint_handler()
    signal.Notify(sig_chan, syscall.SIGINT)

    debugLoop(&d)
}
