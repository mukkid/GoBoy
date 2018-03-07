package GoBoy

import (
	"image"
	"image/color"
)

const (
	VRAMTilePattern      = 0x8000 // 0x8000-0x97FF
	VRAMTilePatternEnd   = 0x97FF // 0x8000-0x97FF
	VRAMBackgroundMap    = 0x9800 // 0x9800-0x9BFF
	VRAMBackgroundMapEnd = 0x9BFF // 0x9800-0x9BFF
)

const (
	TileWidth  = 8  // Each tile is 8 pixels across
	TileHeight = 8  // Each tile is 8 pixels tall
	MapWidth   = 32 // The map is 32 tiles across
	MapHeight  = 32 // The map is 32 tiles tall
)

// paletteMap maps color index to RGBA color. TODO: use register values
var paletteMap = map[uint8]color.RGBA{
	0x00: color.RGBA{0x00, 0x00, 0x00, 0xff},
	0x01: color.RGBA{0x33, 0x33, 0x33, 0xff},
	0x02: color.RGBA{0xcc, 0xcc, 0xcc, 0xff},
	0x03: color.RGBA{0xff, 0xff, 0xff, 0xff},
}

// selectSemiNibble picks the semi-nibble from input given index
func selectSemiNibble(input uint8, index uint8) uint8 {
	return uint8(input>>uint8(index*2)) & 0x03
}

// tileToPixel returns an 8x8 array of RGBA pixel values. It inputs a tile index.
// This is then converted to the raw pointer to the tile data which is a set of
// 8 16-bit unsigned integers. From here it iterates through each 2-bit pixel
// value and translates that into a 32-bit 8x8 array
func tileToPixel(tileIndex uint8, mem *GBMem) [TileHeight][TileWidth]color.RGBA {
	var pixels [TileHeight][TileWidth]color.RGBA
	for y := 0; y < TileHeight; y++ {
		var line [TileWidth]color.RGBA
		for x := 0; x < TileWidth/4; x++ {
			tile := mem.vram[x+(y*TileWidth/4)]
			line[x*4+0] = paletteMap[selectSemiNibble(tile, 0)]
			line[x*4+1] = paletteMap[selectSemiNibble(tile, 1)]
			line[x*4+2] = paletteMap[selectSemiNibble(tile, 2)]
			line[x*4+3] = paletteMap[selectSemiNibble(tile, 3)]
		}
		pixels[y] = line
	}
	return pixels
}

// drawTilePixels draws an 8x8 tile onto an image
func drawTilePixels(image *image.RGBA, pixel [8][8]color.RGBA, xOffset int, yOffset int) *image.RGBA {
	for x := 0; x < TileWidth; x++ {
		for y := 0; y < TileHeight; y++ {
			image.SetRGBA(xOffset*TileWidth+x, yOffset*TileHeight+y, pixel[x][y])
		}
	}
	return image
}

// drawBackground draws the background tile map onto an image
func drawBackground(image *image.RGBA, mem *GBMem) *image.RGBA {
	for x := 0; x < MapWidth; x++ {
		for y := 0; y < MapHeight; y++ {
			// get tile index
			tileIndex := mem.read(uint16(VRAMBackgroundMap + x + (y * MapHeight)))

			// get pixels corresponding to tile index
			pixels := tileToPixel(tileIndex, mem)

			// draw pixels
			drawTilePixels(image, pixels, x, y)
		}
	}
	return image
}
