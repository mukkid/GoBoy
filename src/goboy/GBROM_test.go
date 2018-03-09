package main

import "testing"
import "github.com/stretchr/testify/assert"

/* Test reading from ROM */
func TestGBROMReadROM(t *testing.T) {
	rom := &GBROM{}
	rom.rom[7] = uint8(0xbc)
	assert.Equal(t, uint8(0xbc), rom.readROM(0x7))
}

func TestGBROMLoadROM(t *testing.T) {
	romData := make([]uint8, 0x8000, 0x8000)
	for i, _ := range romData {
		romData[i] = uint8(i % 256)
	}
	gbRom := &GBROM{}
	gbRom.loadROM(romData)
	assert.Equal(t, gbRom.rom[255], uint8(255))
	assert.Equal(t, gbRom.rom[0x7de7], uint8(0xe7))
}

/* Will add tests for readRAM and writeRAM once we've implemented memory faults */
