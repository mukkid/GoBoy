package main

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
	gb.mainMemory.write(address+1, highVal)
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
	gb.mainMemory.write(address+1, highVal)
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

// ADD A, r
func (gb *GameBoy) ADD_a_r(ins [1]uint8) {
	r := Reg8ID(ins[0] & 0x07)
	aVal := gb.get8Reg(A)
	bVal := gb.get8Reg(r)
	out := aVal + bVal
	gb.set8Reg(A, out)
	if (aVal&0x0f)+(bVal&0x0f) > 0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	if (uint16(aVal) + uint16(bVal)) > 0xff {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	gb.regs[PC] += uint16(len(ins))
}

// ADD A, n
func (gb *GameBoy) ADD_a_n(ins [2]uint8) {
	aVal := gb.get8Reg(A)
	bVal := ins[1]
	out := aVal + bVal
	gb.set8Reg(A, out)
	if (aVal&0x0f)+(bVal&0x0f) > 0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if (uint16(aVal) + uint16(bVal)) > 0xff {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// ADD A, (HL)
func (gb *GameBoy) ADD_a_hl(ins [1]uint8) {
	aVal := gb.get8Reg(A)
	bVal := gb.mainMemory.read(gb.get16Reg(HL))
	out := aVal + bVal
	gb.set8Reg(A, out)
	if (aVal&0x0f)+(bVal&0x0f) > 0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if (uint16(aVal) + uint16(bVal)) > 0xff {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// ADC A, r
func (gb *GameBoy) ADC_a_r(ins [1]uint8) {
	r := Reg8ID(ins[0] & 0x07)
	aVal := gb.get8Reg(A)
	bVal := gb.get8Reg(r)
	c := gb.get8Reg(F) & uint8(C_FLAG) >> 4
	out := aVal + bVal + c
	gb.set8Reg(A, out)
	if (aVal&0x0f)+(bVal&0x0f) > 0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if (uint16(aVal) + uint16(bVal)) > 0xff {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// ADC A, n
func (gb *GameBoy) ADC_a_n(ins [2]uint8) {
	aVal := gb.get8Reg(A)
	bVal := ins[1]
	c := gb.get8Reg(F) & uint8(C_FLAG) >> 4
	out := aVal + bVal + c
	gb.set8Reg(A, out)
	if (aVal&0x0f)+(bVal&0x0f) > 0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if (uint16(aVal) + uint16(bVal)) > 0xff {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// ADC A, (HL)
func (gb *GameBoy) ADC_a_hl(ins [1]uint8) {
	aVal := gb.get8Reg(A)
	bVal := gb.mainMemory.read(gb.get16Reg(HL))
	c := gb.get8Reg(F) & uint8(C_FLAG) >> 4
	out := aVal + bVal + c
	gb.set8Reg(A, out)
	if (aVal&0x0f)+(bVal&0x0f) > 0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if (uint16(aVal) + uint16(bVal)) > 0xff {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// SUB A, r
func (gb *GameBoy) SUB_a_r(ins [1]uint8) {
	r := Reg8ID(ins[0] & 0x07)
	aVal := gb.get8Reg(A)
	bVal := gb.get8Reg(r)
	out := aVal - bVal
	gb.set8Reg(A, out)
	if aVal < bVal {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	if aVal&0x0f < bVal&0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// SUB A, n
func (gb *GameBoy) SUB_a_n(ins [2]uint8) {
	aVal := gb.get8Reg(A)
	bVal := ins[1]
	out := aVal - bVal
	gb.set8Reg(A, out)
	if aVal < bVal {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	if aVal&0x0f < bVal&0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// SUB A, (HL)
func (gb *GameBoy) SUB_a_hl(ins [1]uint8) {
	aVal := gb.get8Reg(A)
	bVal := gb.mainMemory.read(gb.get16Reg(HL))
	out := aVal - bVal
	gb.set8Reg(A, out)
	if aVal < bVal {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	if aVal&0x0f < bVal&0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// SBC A, r
func (gb *GameBoy) SBC_a_r(ins [1]uint8) {
	r := Reg8ID(ins[0] & 0x07)
	aVal := gb.get8Reg(A)
	bVal := gb.get8Reg(r)
	c := gb.get8Reg(F) & uint8(C_FLAG) >> 4
	out := aVal - (bVal + c)
	gb.set8Reg(A, out)
	if aVal < (bVal + c) {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	if aVal&0x0f < (bVal+c)&0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// SBC A, n
func (gb *GameBoy) SBC_a_n(ins [2]uint8) {
	aVal := gb.get8Reg(A)
	bVal := ins[1]
	c := gb.get8Reg(F) & uint8(C_FLAG) >> 4
	out := aVal - (bVal + c)
	gb.set8Reg(A, out)
	if aVal < (bVal + c) {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	if aVal&0x0f < (bVal+c)&0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// SBC A, (HL)
func (gb *GameBoy) SBC_a_hl(ins [1]uint8) {
	aVal := gb.get8Reg(A)
	bVal := gb.mainMemory.read(gb.get16Reg(HL))
	c := gb.get8Reg(F) & uint8(C_FLAG) >> 4
	out := aVal - (bVal + c)
	gb.set8Reg(A, out)
	if aVal < (bVal + c) {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	if aVal&0x0f < (bVal+c)&0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// AND A, r
func (gb *GameBoy) AND_a_r(ins [1]uint8) {
	r := Reg8ID(ins[0] & 0x07)
	aVal := gb.get8Reg(A)
	bVal := gb.get8Reg(r)
	out := aVal & bVal
	gb.set8Reg(A, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.modifyFlag(H_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// AND A, n
func (gb *GameBoy) AND_a_n(ins [2]uint8) {
	aVal := gb.get8Reg(A)
	bVal := ins[1]
	out := aVal & bVal
	gb.set8Reg(A, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.modifyFlag(H_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// AND A, (HL)
func (gb *GameBoy) AND_a_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	aVal := gb.get8Reg(A)
	bVal := gb.mainMemory.read(address)
	out := aVal & bVal
	gb.set8Reg(A, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.modifyFlag(H_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// OR A, r
func (gb *GameBoy) OR_a_r(ins [1]uint8) {
	r := Reg8ID(ins[0] & 0x07)
	aVal := gb.get8Reg(A)
	bVal := gb.get8Reg(r)
	out := aVal | bVal
	gb.set8Reg(A, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// OR A, n
func (gb *GameBoy) OR_a_n(ins [2]uint8) {
	aVal := gb.get8Reg(A)
	bVal := ins[1]
	out := aVal | bVal
	gb.set8Reg(A, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// OR A, (HL)
func (gb *GameBoy) OR_a_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	aVal := gb.get8Reg(A)
	bVal := gb.mainMemory.read(address)
	out := aVal | bVal
	gb.set8Reg(A, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// XOR A, r
func (gb *GameBoy) XOR_a_r(ins [1]uint8) {
	r := Reg8ID(ins[0] & 0x07)
	aVal := gb.get8Reg(A)
	bVal := gb.get8Reg(r)
	out := aVal ^ bVal
	gb.set8Reg(A, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// XOR A, n
func (gb *GameBoy) XOR_a_n(ins [2]uint8) {
	aVal := gb.get8Reg(A)
	bVal := ins[1]
	out := aVal ^ bVal
	gb.set8Reg(A, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// XOR A, (HL)
func (gb *GameBoy) XOR_a_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	aVal := gb.get8Reg(A)
	bVal := gb.mainMemory.read(address)
	out := aVal ^ bVal
	gb.set8Reg(A, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// CP A, r
func (gb *GameBoy) CP_a_r(ins [1]uint8) {
	r := Reg8ID(ins[0] & 0x07)
	aVal := gb.get8Reg(A)
	bVal := gb.get8Reg(r)
	if aVal > bVal {
		gb.modifyFlag(Z_FLAG, CLEAR)
		gb.modifyFlag(H_FLAG, SET)
		gb.modifyFlag(C_FLAG, CLEAR)
	} else if aVal < bVal {
		gb.modifyFlag(Z_FLAG, CLEAR)
		gb.modifyFlag(H_FLAG, CLEAR)
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, SET)
		gb.modifyFlag(H_FLAG, CLEAR)
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// CP A, n
func (gb *GameBoy) CP_a_n(ins [2]uint8) {
	aVal := gb.get8Reg(A)
	bVal := ins[1]
	if aVal > bVal {
		gb.modifyFlag(Z_FLAG, CLEAR)
		gb.modifyFlag(H_FLAG, SET)
		gb.modifyFlag(C_FLAG, CLEAR)
	} else if aVal < bVal {
		gb.modifyFlag(Z_FLAG, CLEAR)
		gb.modifyFlag(H_FLAG, CLEAR)
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, SET)
		gb.modifyFlag(H_FLAG, CLEAR)
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// CP A,(HL)
func (gb *GameBoy) CP_a_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	aVal := gb.get8Reg(A)
	bVal := gb.mainMemory.read(address)
	if aVal > bVal {
		gb.modifyFlag(Z_FLAG, CLEAR)
		gb.modifyFlag(H_FLAG, SET)
		gb.modifyFlag(C_FLAG, CLEAR)
	} else if aVal < bVal {
		gb.modifyFlag(Z_FLAG, CLEAR)
		gb.modifyFlag(H_FLAG, CLEAR)
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, SET)
		gb.modifyFlag(H_FLAG, CLEAR)
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// INC r
func (gb *GameBoy) INC_r(ins [1]uint8) {
	r := Reg8ID((ins[0] >> 3) & 0x07)
	val := gb.get8Reg(r)
	out := val + 1
	gb.set8Reg(r, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	if (val&0x0f)+0x01 > 0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// INC (HL)
func (gb *GameBoy) INC_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	val := gb.mainMemory.read(address)
	out := val + 1
	gb.mainMemory.write(address, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	if (val&0x0f)+0x01 > 0x0f {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// DEC r
func (gb *GameBoy) DEC_r(ins [1]uint8) {
	r := Reg8ID((ins[0] >> 3) & 0x07)
	val := gb.get8Reg(r)
	out := val - 1
	gb.set8Reg(r, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	if val&0x0f < 0x01 {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

// DEC (HL)
func (gb *GameBoy) DEC_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	val := gb.mainMemory.read(address)
	out := val - 1
	gb.mainMemory.write(address, out)
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	if val&0x0f < 0x01 {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	gb.modifyFlag(N_FLAG, SET)
	gb.regs[PC] += uint16(len(ins))
}

/* 16 BIT ALU OPCODES */

// ADD HL ss
func (gb *GameBoy) ADD_hl_ss(ins [1]uint8) {
	aVal := gb.get16Reg(HL)
	ss := Reg16ID((ins[0] >> 4) & 0x03)
	bVal := gb.get16Reg(ss)
	out := aVal + bVal
	gb.set16Reg(HL, out)
	gb.modifyFlag(N_FLAG, CLEAR)
	if (aVal&0x0fff)+(bVal&0x0fff) > 0x0fff {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if uint32(aVal)+uint32(bVal) > 0xffff {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	gb.regs[PC] += uint16(len(ins))
}

// ADD SP e
func (gb *GameBoy) ADD_sp_e(ins [2]uint8) {
	aVal := gb.get16Reg(SP)
	bVal := uint16(ins[1])
	out := aVal + bVal
	gb.set16Reg(SP, out)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.modifyFlag(Z_FLAG, CLEAR)
	if (aVal&0x0fff)+(bVal&0x0fff) > 0x0fff {
		gb.modifyFlag(H_FLAG, SET)
	} else {
		gb.modifyFlag(H_FLAG, CLEAR)
	}
	if uint32(aVal)+uint32(bVal) > 0xffff {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}
	gb.regs[PC] += uint16(len(ins))
}

// TODO: Implement 16 bit INC
// TODO: Implement 16 bit DEC

func (gb *GameBoy) JP_nn(ins [3]uint8) {
	address := binary.LittleEndian.Uint16(ins[1:])
	gb.set16Reg(PC, address)
	gb.regs[PC] += uint16(len(ins))
}
