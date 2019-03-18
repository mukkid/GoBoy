package main

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

type color struct {
	r, g, b, a byte
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*visibleWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index] = c.g
		pixels[index] = c.b
		pixels[index] = c.a
	}
}

// paletteMap maps color index to RGBA color. TODO: use register values
var paletteMap = map[uint8]color{
	0x00: color{0x00, 0x00, 0x00, 0xff},
	0x01: color{0x33, 0x33, 0x33, 0xff},
	0x02: color{0xcc, 0xcc, 0xcc, 0xff},
	0x03: color{0xff, 0xff, 0xff, 0xff},
}

// selectSemiNibble picks the semi-nibble from input given index
func selectSemiNibble(input uint8, index uint8) uint8 {
	return uint8(input>>uint8(index*2)) & 0x03
}

func compositePixel(lsb, hsb, index uint8) uint8 {
	return (lsb>>index)&0x1 | (hsb>>index)&0x1<<1
}

// tileToPixel returns an 8x8 array of RGBA pixel values. It inputs a tile index.
// This is then converted to the raw pointer to the tile data which is a set of
// 8 16-bit unsigned integers. From here it iterates through each 2-bit pixel
// value and translates that into a 32-bit 8x8 array
func tileToPixel(tileIndex uint8, mem *GBMem) []byte {
	// var pixels [TileHeight * TileWidth * 4]byte
	pixels := make([]byte, 64)
	for y := 0; y < TileHeight; y++ {
		lsb := mem.vram[int(tileIndex)*16+2*y]
		hsb := mem.vram[int(tileIndex)*16+2*y+1]

		for x := 0; x < TileWidth; x++ {
			// pixels[(y*TileHeight)+(TileWidth-1-x)] = paletteMap[compositePixel(lsb, hsb, uint8(x))]
			mappedColor := paletteMap[compositePixel(lsb, hsb, uint8(x))]
			setPixel(x, y, mappedColor, pixels)
		}
	}
	return pixels
}

func updateSlice(largerSlice []byte, smallerSlice []byte, offset int) {
	for i := 0; i < len(smallerSlice); i++ {
		largerSlice[offset+i] = smallerSlice[i]
	}
}

// drawBackground draws the background tile map onto an image
func drawBackground(lcd []byte, mem *GBMem) []byte {
	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {
			// get tile index
			tileIndex := mem.read(uint16(VRAMBackgroundMap + x + (y * MapHeight)))

			// get pixels corresponding to tile index
			pixels := tileToPixel(tileIndex, mem)

			// draw pixels
			updateSlice(lcd, pixels, x+y*MapHeight)
		}
	}
	return lcd
}
