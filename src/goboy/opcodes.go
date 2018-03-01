package GoBoy

// Load r1 <- r2
func (gb *GameBoy) LD_r_r(ins uint8) {
    r1 := Reg8ID((ins & 0x38) >> 3)
    r2 := Reg8ID(ins & 0x07)
    gb.set8Reg(r1, gb.get8Reg(r2))
    gb.regs[PC] += 1
}

// Load r <- n
func (gb *GameBoy) LD_r_n(ins uint16) {
    n := uint8(ins)
    r := Reg8ID((ins & 0x3800) >> 11)
    gb.set8Reg(r, n)
    gb.regs[PC] += 2
}

// Load r <- (HL)
func (gb *GameBoy) LD_r_hl(ins uint8) {
    r := Reg8ID((ins & 0x38) >> 3)
    address := gb.get16Reg(HL)
    gb.set8Reg(r, gb.mainMemory.read(address))
    gb.regs[PC] += 1
}


/*
 * TODO: Write unit test
 * NOTE: gb.mainMemory.write() has not been implemented yet
 * Load (HL) <- r
 */
func (gb *GameBoy) LD_hl_r(ins uint8) {
    r := Reg8ID(ins & 0x07)
    address := gb.get16Reg(HL)
    gb.mainMemory.write(address, gb.get8Reg(r))
    gb.regs[PC] += 1
}

/* 
 * TODO: Write unit test
 * NOTE: gb.mainMemory.write() has not been implemented yet
 * Load (HL) <- n
 */
func (gb *GameBoy) LD_hl_n(ins uint16) {
    n := uint8(ins)
    address := gb.get16Reg(HL)
    gb.mainMemory.write(address, n)
    gb.regs[PC] += 2
}

/*
 * Load A <- (BC)
 */
func (gb *GameBoy) LD_a_bc(ins uint8) {
    address := gb.get16Reg(BC)
    gb.set8Reg(A, gb.mainMemory.read(address))
    gb.regs[PC] += 1
}

/*
 * Load A <- (DE)
 */
func (gb *GameBoy) LD_a_de(ins uint8) {
    address := gb.get16Reg(DE)
    gb.set8Reg(A, gb.mainMemory.read(address))
    gb.regs[PC] += 1
}

/*
 * TODO: Write unit test
 * Load A <- (nn)
 * NOTE: By convention, 24 bit instructions will be passed
 * as uint32 numbers with padding on the MSBs
 */
func (gb *GameBoy) LD_a_nn(ins uint32) {
    addressFlipped := uint16(ins)
    addressLow := addressFlipped >> 8
    addressHigh := addressFlipped << 8
    address := addressLow | addressHigh
    gb.set8Reg(A, gb.mainMemory.read(address))
    gb.regs[PC] += 3
}

/*
 * TODO: Write unit test
 * NOTE: gb.mainMemory.write() needs to be implemented
 * Load (BC) <- A
 */
func (gb *GameBoy) LD_bc_a(ins uint8) {
    address := gb.get16Reg(BC)
    gb.mainMemory.write(address, gb.get8Reg(A))
    gb.regs[PC] += 1
}

/*
 * TODO: Write unit test
 * NOTE: gb.mainMemory.write() needs to be implemented
 * Load (DE) <- A
 */
func (gb *GameBoy) LD_de_a(ins uint8) {
    address := gb.get16Reg(DE)
    gb.mainMemory.write(address, gb.get8Reg(A))
    gb.regs[PC] += 1
}

/*
 *NOTE: gb.mainMemory.write() needs to be implemented
 * Load (nn) <- A
 */
func (gb *GameBoy) LD_nn_a(ins uint32) {
    addressFlipped := uint16(ins)
    addressLow := addressFlipped >> 8
    addressHigh := addressFlipped << 8
    address := addressLow | addressHigh
    gb.mainMemory.write(address, gb.get8Reg(A))
    gb.regs[PC] += 3
}
