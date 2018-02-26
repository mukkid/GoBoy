package GoBoy

// Load r1 <- r2
func (gb *GameBoy) LD_r_r(ins uint8) {
    r1 := Reg8ID((ins & 0x38) >> 3)
    r2 := Reg8ID(ins & 0x07)
    gb.set8Reg(r1, gb.get8Reg(r2))
    gb.incrementPc(1)
}

// Load r <- n
func (gb *GameBoy) LD_r_n(ins uint16) {
    n := uint8(ins)
    r := Reg8ID((ins & 0x3800) >> 11)
    gb.set8Reg(r, n)
    gb.incrementPc(2)
}
