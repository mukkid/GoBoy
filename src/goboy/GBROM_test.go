package GoBoy

import "testing"
import "github.com/stretchr/testify/assert"

/* Test reading from ROM */
func TestGBROMReadROM(t *testing.T) {
    rom := &GBROM{}
    rom.rom1[7] = uint8(0xbc);
    assert.Equal(t, uint8(0xbc), rom.readROM(0x4007))
}

/* Will add tests for readRAM and writeRAM once we've implemented memory faults */
