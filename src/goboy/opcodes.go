package GoBoy

import "encoding/binary"

// Load r1 <- r2
func (gb *GameBoy) LD_r_r(ins [1]uint8) {
    r1 := Reg8ID((ins[0] & 0x38) >> 3)
    r2 := Reg8ID(ins[0] & 0x07)
    gb.set8Reg(r1, gb.get8Reg(r2))
    gb.regs[PC] += uint16(len(ins))
}

// Load r <- n
func (gb *GameBoy) LD_r_n(ins [2]uint8) {
    n := ins[1]
    r := Reg8ID(ins[0] >> 3)
    gb.set8Reg(r, n)
    gb.regs[PC] += uint16(len(ins))
}

// Load r <- (HL)
func (gb *GameBoy) LD_r_hl(ins [1]uint8) {
    r := Reg8ID((ins[0] & 0x38) >> 3)
    address := gb.get16Reg(HL)
    gb.set8Reg(r, gb.mainMemory.read(address))
    gb.regs[PC] += uint16(len(ins))
}


// Load (HL) <- r
func (gb *GameBoy) LD_hl_r(ins [1]uint8) {
    r := Reg8ID(ins[0] & 0x07)
    address := gb.get16Reg(HL)
    gb.mainMemory.write(address, gb.get8Reg(r))
    gb.regs[PC] += uint16(len(ins))
}

// Load (HL) <- n
func (gb *GameBoy) LD_hl_n(ins [2]uint8) {
    n := ins[1]
    address := gb.get16Reg(HL)
    gb.mainMemory.write(address, n)
    gb.regs[PC] += uint16(len(ins))
}

// Load A <- (BC)
func (gb *GameBoy) LD_a_bc(ins [1]uint8) {
    address := gb.get16Reg(BC)
    gb.set8Reg(A, gb.mainMemory.read(address))
    gb.regs[PC] += uint16(len(ins))
}

// Load A <- (DE)
func (gb *GameBoy) LD_a_de(ins [1]uint8) {
    address := gb.get16Reg(DE)
    gb.set8Reg(A, gb.mainMemory.read(address))
    gb.regs[PC] += uint16(len(ins))
}

/*
 * Load A <- (nn)
 * NOTE: By convention, 24 bit instructions will be passed
 * as uint32 numbers with padding on the MSBs
 */
func (gb *GameBoy) LD_a_nn(ins [3]uint8) {
    address := binary.LittleEndian.Uint16(ins[1:])
    gb.set8Reg(A, gb.mainMemory.read(address))
    gb.regs[PC] += uint16(len(ins))
}

// Load (BC) <- A
func (gb *GameBoy) LD_bc_a(ins [1]uint8) {
    address := gb.get16Reg(BC)
    gb.mainMemory.write(address, gb.get8Reg(A))
    gb.regs[PC] += uint16(len(ins))
}

// Load (DE) <- A
func (gb *GameBoy) LD_de_a(ins [1]uint8) {
    address := gb.get16Reg(DE)
    gb.mainMemory.write(address, gb.get8Reg(A))
    gb.regs[PC] += uint16(len(ins))
}

// Load (nn) <- A
func (gb *GameBoy) LD_nn_a(ins [3]uint8) {
    address := binary.LittleEndian.Uint16(ins[1:])
    gb.mainMemory.write(address, gb.get8Reg(A))
    gb.regs[PC] += uint16(len(ins))
}

/* 16 BIT LOADS */

// LD dd <- nn
func (gb *GameBoy) LD_dd_nn(ins [3]uint8) {
    reg := Reg16ID(ins[0] >> 4)
    nn := binary.LittleEndian.Uint16(ins[1:])
    gb.set16Reg(reg, nn)
    gb.regs[PC] += uint16(len(ins))
}

// LD HL <- (nn)
func (gb *GameBoy) LD_hl_nn(ins [3]uint8) {
    address := binary.LittleEndian.Uint16(ins[1:])
    lowVal := gb.mainMemory.read(address)
    highVal := gb.mainMemory.read(address + 1)
    gb.set8Reg(H, highVal)
    gb.set8Reg(L, lowVal)
    gb.regs[PC] += uint16(len(ins))
}

// LD dd <- (NN)
func (gb *GameBoy) LD_dd_NN(ins [4]uint8) {
    reg := Reg16ID((ins[1] >> 4) & 0x03)
    address := binary.LittleEndian.Uint16(ins[2:])
    lowVal := gb.mainMemory.read(address)
    highVal := gb.mainMemory.read(address + 1)
    val := uint16(lowVal) | (uint16(highVal) << 8)
    gb.set16Reg(reg, val)
    gb.regs[PC] += uint16(len(ins))
}

// LD (nn) <- HL
func (gb *GameBoy) LD_nn_hl(ins [3]uint8) {
    address := binary.LittleEndian.Uint16(ins[1:])
    highVal := gb.get8Reg(H)
    lowVal := gb.get8Reg(L)
    gb.mainMemory.write(address, lowVal)
    gb.mainMemory.write(address + 1, highVal)
    gb.regs[PC] += uint16(len(ins))
}

// LD (nn) <- dd
func (gb *GameBoy) LD_nn_dd(ins [4]uint8) {
    reg := Reg16ID((ins[1] >> 4) & 0x03)
    value := gb.get16Reg(reg)
    lowVal := uint8(value)
    highVal := uint8(value >> 8)
    address := binary.LittleEndian.Uint16(ins[2:])
    gb.mainMemory.write(address, lowVal)
    gb.mainMemory.write(address + 1, highVal)
    gb.regs[PC] += uint16(len(ins))
}

// LD SP <- HL
func (gb *GameBoy) LD_sp_hl(ins [1]uint8) {
    gb.set16Reg(SP, gb.get16Reg(HL))
    gb.regs[PC] += uint16(len(ins))
}

// PUSH qq
func (gb *GameBoy) PUSH_qq(ins [1]uint8) {
    reg := Reg16ID((ins[0] >> 4) & 0x3)
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
    gb.regs[PC] += uint16(len(ins))
}

// POP qq
func (gb *GameBoy) POP_qq(ins [1]uint8) {
    reg := Reg16ID((ins[0] >> 4) & 0x3)
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
