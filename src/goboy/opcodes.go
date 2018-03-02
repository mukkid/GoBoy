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
 * TODO: Write unit test
 * NOTE: gb.mainMemory.write() needs to be implemented
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

/* 16 BIT LOADS */

// LD dd <- nn
func (gb *GameBoy) LD_dd_nn(ins uint32) {
    reg := Reg16ID(ins >> 20)
    nn := uint16(ins)
    gb.set16Reg(reg, nn)
    gb.regs[PC] += 3
}

// LD HL <- (nn)
func (gb *GameBoy) LD_hl_nn(ins uint32) {
    address := uint16(ins)
    lowVal := gb.mainMemory.read(address)
    highVal := gb.mainMemory.read(address + 1)
    gb.set8Reg(H, highVal)
    gb.set8Reg(L, lowVal)
    gb.regs[PC] += 3
}

// LD dd <- (NN)
func (gb *GameBoy) LD_dd_NN(ins uint32) {
    reg := Reg16ID((ins >> 20) & 0x03)
    address := uint16(ins)
    lowVal := gb.mainMemory.read(address)
    highVal := gb.mainMemory.read(address + 1)
    val := uint16(lowVal) | (uint16(highVal) << 8)
    gb.set16Reg(reg, val)
    gb.regs[PC] += 4
}

// TODO: Write unit test
// NOTE: gb.mainMemory.write() needs to be implemented
// LD (nn) <- HL
func (gb *GameBoy) LD_nn_hl(ins uint32) {
    address := uint16(ins)
    highVal := gb.get8Reg(H)
    lowVal := gb.get8Reg(L)
    gb.mainMemory.write(address, lowVal)
    gb.mainMemory.write(address + 1, highVal)
    gb.regs[PC] += 3
}

// TODO: Write unit test
// NOTE: gb.mainMemory.write() needs to be implemented
// LD (nn) <- dd
func (gb *GameBoy) LD_nn_dd(ins uint32) {
    reg := Reg16ID((ins >> 20) & 0x03)
    value := gb.get16Reg(reg)
    lowVal := uint8(value)
    highVal := uint8(value >> 8)
    address := uint16(ins)
    gb.mainMemory.write(address, lowVal)
    gb.mainMemory.write(address + 1, highVal)
    gb.regs[PC] += 4
}

// LD SP <- HL
func (gb *GameBoy) LD_sp_hl(ins uint8) {
    gb.set16Reg(SP, gb.get16Reg(HL))
    gb.regs[PC] += 1
}

// TODO: Write unit test
// NOTE: gb.mainMemory.write() needs to be implemented
// PUSH qq
func (gb *GameBoy) PUSH_qq(ins uint8) {
    reg := Reg16ID((ins >> 4) & 0x3)
    if reg == 0x3 {
        reg = 0x5
    }
    val := gb.get16Reg(reg)
    lowVal := uint8(val)
    highVal := uint8(val >> 8)
    gb.regs[SP] -= 1
    gb.mainMemory.write(gb.get16Reg(SP), highVal)
    gb.regs[SP] -= 1
    gb.mainMemory.write(gb.get16Reg(SP), lowVal)
    gb.regs[PC] += 1
}

// TODO: Write unit test
// NOTE: Accessing RAM needs to be implemented
// POP qq
func (gb *GameBoy) POP_qq(ins uint8) {
    reg := Reg16ID((ins >> 4) & 0x3)
    if reg == 0x3 {
        reg = 0x5
    }
    lowVal := gb.mainMemory.read(gb.get16Reg(SP))
    gb.regs[SP] += 1
    highVal := gb.mainMemory.read(gb.get16Reg(SP))
    gb.regs[SP] += 1
    val := uint16(lowVal) | (uint16(highVal) << 8)
    gb.set16Reg(reg, val)
    gb.regs[PC] += 1
}
