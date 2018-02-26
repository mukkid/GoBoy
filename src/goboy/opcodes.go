package GoBoy

// Load r1 <- r2
func (gb *GameBoy) LD_r_r(ins uint8) {
    r2 := Reg8ID(ins & 0x07)
    r1 := Reg8ID((ins & 0x38) >> 3)
    gb.set8Reg(r1, gb.get8Reg(r2))
    gb.incrementPc(1)
}
