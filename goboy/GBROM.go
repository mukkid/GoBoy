package main

import (
	"bytes"
	"io/ioutil"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/*
 * Implements the ROM only cartridge, the simplest one
 */
type GBROM struct {
	/* "Switchable" ROM1 bank. It's not switchable. */
	rom [0x8000]uint8
}

func (r *GBROM) readROM(addr uint16) uint8 {
	/* ROM1 bank starts at 0x4000 in the main memory addr space */
	return r.rom[addr]
}

func (r *GBROM) readRAM(addr uint16) uint8 {
	/*
	 * ROM-only cartrdiges have no RAM
	 * TODO: fault
	 */
	return 0
}

func (r *GBROM) writeROM(addr uint16, data uint8) {
	/*
	 * Not sure if this cartridge type crashes on ROM writes
	 */
	return
}

func (r *GBROM) writeRAM(addr uint16, data uint8) {
	/*
	 * ROM-only cartrdiges have no RAM
	 * TODO: fault
	 */
	return
}

func (r *GBROM) loadROM(data []uint8) error {
	copy(r.rom[:], data)
	return nil
}

func (r *GBROM) loadROMFromFile(fname string) error {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		println(err)
		return err
	}
	copy(r.rom[:], data)
	return err
}

func (r *GBROM) reader() *bytes.Reader {
	return bytes.NewReader(r.rom[:])
}

func newGBROM() *GBROM {
	return &GBROM{}
}
