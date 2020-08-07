package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"time"
	"image"
	//"os"
)

const (
	GUIPeriod = 10 * time.Microsecond
)

type Gooey struct {
	window *sdl.Window
	renderer *sdl.Renderer
	texture *sdl.Texture
	rect sdl.Rect
	/* Clock for the GUI */
	clock *time.Ticker
	/* Only one thread is writing to quit, so no synchronization required */
	quit bool
	debugger *Debugger
}

func newGooey(d *Debugger) (*Gooey, error) {
	g := Gooey{
		quit  	: false,
		clock 	: time.NewTicker(GUIPeriod),
		debugger : d,
	}
	var err error
	g.window, err = sdl.CreateWindow(
		"Goboy",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		640,
		480,
		sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}
	g.renderer, err = sdl.CreateRenderer(g.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}
	//g.renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
	g.texture, err = g.renderer.CreateTexture(
		sdl.PIXELFORMAT_ARGB8888,
		sdl.TEXTUREACCESS_STREAMING,
		screenWidth,
		screenHeight)
	if err != nil {
		return nil, err
	}

	g.rect = sdl.Rect{
		X: 0,
		Y: 0,
		W: screenWidth,
		H: screenHeight,
	}
	return &g, nil
}

func (g *Gooey) lockTexture() ([]byte, int, error) {
	return g.texture.Lock(&g.rect)
}

func (g *Gooey) unlockTexture() () {
	g.texture.Unlock()
}

/* Update the window with the new texture pixels */
func (g *Gooey) renderUpdate() error {
	err := g.renderer.Copy(g.texture, &g.rect, &g.rect)
	if err != nil {
		return err
	}
	g.renderer.Present()
	return nil
}

func (g *Gooey) eventLoop() error {
	bgImage := image.NewRGBA(image.Rect(0, 0, 256, 256))
	for _ = range g.clock.C {
		if g.debugger.gb.Halt {
			break
		}
		/* Render stuff */
		pixels, _, err := g.lockTexture()
		if err != nil {
			fmt.Println("oh no")
			return err
		}

		drawBackground(bgImage, g.debugger.gb.mainMemory)
		copy(pixels, bgImage.Pix)

		//file, _ := os.OpenFile("pixels", os.O_RDWR|os.O_CREATE, 0755)
		//file.Write(pixels)
		//file.Close()

		g.unlockTexture()
		err = g.renderUpdate()
		if err != nil {
			fmt.Println("oh no")
			return err
		}

		/* Handle events */
		event := sdl.PollEvent()
		if event != nil && event.GetType() == sdl.QUIT {
			g.quit = true
			break
		}
	}
	sdl.Quit()
	return nil
}
