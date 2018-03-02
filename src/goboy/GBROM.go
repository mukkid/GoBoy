package GoBoy

/*
 * Implements the ROM only cartridge, the simplest one
 */
type GBROM struct {
    /* "Switchable" ROM1 bank. It's not switchable. */
    rom1 [16 * 1024] uint8
}

func (r *GBROM) readROM(addr uint16) uint8 {
    /* ROM1 bank starts at 0x4000 in the main memory addr space */
    return r.rom1[addr - 0x4000];
}

func (r *GBROM) readRAM(addr uint16) uint8 {
    /*
     * ROM-only cartrdiges have no RAM
     * TODO: fault
     */
    return 0;
}

func (r *GBROM) writeROM(addr uint16, data uint8) {
    /*
     * Not sure if this cartridge type crashes on ROM writes
     */
     return;
}

func (r *GBROM) writeRAM(addr uint16, data uint8) {
    /*
     * ROM-only cartrdiges have no RAM
     * TODO: fault
     */
    return;
}
