package main

import (
	"flag"
	"fmt"
	"image"
	"os/signal"
	"syscall"
	"time"
	"github.com/veandco/go-sdl2/sdl"
)

/* 256x256 is written to in total but only 160x144 is visible. */
const (
	screenWidth   = 256
	screenHeight  = 256
	visibleWidth  = 160
	visibleHeight = 144
)

func main() {
	/* Initialize SDL */
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return
	}
	/* init gameboy */
	Gb := &GameBoy{
		Register:         &Register{},
		mainMemory:       &GBMem{cartridge: &GBROM{}},
		interruptEnabled: true,
		image:            image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
		LCDClock:         time.NewTicker(108 * time.Microsecond),
		CPUClock:         time.NewTicker(GBClockPeriod),
		TSC:              0,
		TSCStart:         0,
		Paused:           true,
	}

	/* load rom from file */
	rom_path := flag.String("rom", "", "rom image to load")
	flag.Parse()
	if *rom_path != "" {
		Gb.mainMemory.cartridge.loadROMFromFile(*rom_path)
		fmt.Printf("Loaded %s\n", *rom_path)
	}

	/* Initialize joypad values */
	Gb.mainMemory.ioregs[0] = 0xff

	/* Initialize PC to 0x100 */
	Gb.set16Reg(PC, 0x100)

	d := NewDebugger(Gb)
	/* Initialize SIGINT handler */
	go d.SIGINTListener()
	signal.Notify(sig_chan, syscall.SIGINT)

	go Gb.LCDLoop()
	go Gb.TSCLoop()
	go d.debugLoop()
	g, err := newGooey(d)
	if err != nil {
		return
	}
	/* Main SDL loop */
	err = g.eventLoop()
}
