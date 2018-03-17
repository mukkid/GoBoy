package main

import "testing"
import "github.com/stretchr/testify/assert"

/* Test reading working RAM */
func TestGBMemReadRAM(t *testing.T) {
	mem := &GBMem{}
	addr := uint16(0xc070)
	mem.wram[addr-0xc000] = uint8(0xfa)
	assert.Equal(t, mem.read(addr), uint8(0xfa))
}

/* Test reading echo of 0xe000 - 0xfe00 */
func TestGBMemReadRAMEcho(t *testing.T) {
	mem := &GBMem{}
	addr := uint16(0xc070)
	mem.wram[addr-0xc000] = uint8(0xfa)
	addr += 0x2000
	assert.Equal(t, mem.read(addr), uint8(0xfa))
}

/* Test reading HRAM */
func TestGBMemReadHRAM(t *testing.T) {
	mem := &GBMem{}
	addr := uint16(0xff85)
	mem.hram[addr-0xff80] = 0xcb
	assert.Equal(t, mem.read(addr), uint8(0xcb))
}

/* Test writing working RAM */
func TestGBMemWriteRAM(t *testing.T) {
	mem := &GBMem{}
	addr := uint16(0xc0c0)
	mem.write(addr, uint8(0xfa))
	assert.Equal(t, mem.wram[addr-0xc000], uint8(0xfa))
}

/* Test reading echo of 0xe000 - 0xfe00 */
func TestGBMemWriteRAMEcho(t *testing.T) {
	mem := &GBMem{}
	addr := uint16(0xf080)
	mem.write(addr, 0x75)
	addr -= 0x2000
	assert.Equal(t, mem.wram[addr-0xc000], uint8(0x75))
}

/* Test reading HRAM */
func TestGBMemWriteHRAM(t *testing.T) {
	mem := &GBMem{}
	addr := uint16(0xfffe)
	mem.write(addr, 0x1c)
	assert.Equal(t, mem.hram[addr-0xff80], uint8(0x1c))
}

/* Test reading n bytes */
func TestReadN(t *testing.T) {
	mem := &GBMem{}
	mem.write(0xff85, 0x12)
	mem.write(0xff86, 0x34)
	mem.write(0xff87, 0x56)
	assert.Equal(t, mem.readN(0xff85, 3), []uint8{0x12, 0x34, 0x56})
	assert.Equal(t, mem.readN(0xff85, 2), []uint8{0x12, 0x34})
	assert.Equal(t, mem.readN(0xff85, 1), []uint8{0x12})
}
