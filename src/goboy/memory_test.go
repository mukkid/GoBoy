package GoBoy

import "testing"
import "github.com/stretchr/testify/assert"

/* Test reading ROM bank 0 */
func TestGBMemReadROM(t *testing.T) {
    mem := &GBMem{}
    addr := uint16(0x0070)
    mem.rom0[addr] = uint8(0xfa)
    assert.Equal(t, mem.read(addr), uint8(0xfa))
}

/* Test reading working RAM */
func TestGBMemReadRAM(t *testing.T) {
    mem := &GBMem{}
    addr := uint16(0xc070)
    mem.wram[addr - 0xc000 ] = uint8(0xfa)
    assert.Equal(t, mem.read(addr), uint8(0xfa))
}

/* Test reading echo of 0xe000 - 0xfe00 */
func TestGBMemReadRAMEcho(t *testing.T) {
    mem := &GBMem{}
    addr := uint16(0xc070)
    mem.wram[addr - 0xc000 ] = uint8(0xfa)
    addr += 0x2000
    assert.Equal(t, mem.read(addr), uint8(0xfa))
}

/* Test reading HRAM */
func TestGBMemReadHRAM(t *testing.T) {
    mem := &GBMem{}
    addr := uint16(0xff85)
    mem.hram[addr - 0xff80] = 0xcb
    assert.Equal(t, mem.read(addr), uint8(0xcb))
}

/* Test writing working RAM */
func TestGBMemWriteRAM(t *testing.T) {
    mem := &GBMem{}
    addr := uint16(0xc0c0)
    mem.write(addr, uint8(0xfa))
    assert.Equal(t, mem.wram[addr - 0xc000], uint8(0xfa))
}

/* Test reading echo of 0xe000 - 0xfe00 */
func TestGBMemWriteRAMEcho(t *testing.T) {
    mem := &GBMem{}
    addr := uint16(0xf080)
    mem.write(addr, 0x75)
    addr -= 0x2000
    assert.Equal(t, mem.wram[addr - 0xc000], uint8(0x75))
}

/* Test reading HRAM */
func TestGBMemWriteHRAM(t *testing.T) {
    mem := &GBMem{}
    addr := uint16(0xfffe)
    mem.write(addr, 0x1c)
    assert.Equal(t, mem.hram[addr - 0xff80], uint8(0x1c))
}
