package main

import (
	"encoding/binary"
	"math/bits"
)

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
}

func (gb *GameBoy) JP_cc_nn(ins [3]uint8) {
	Z := gb.getFlag(Z_FLAG)
	C := gb.getFlag(C_FLAG)
	cc := ins[0] >> 3 & 0x03
	if cc == 0x00 && Z == 0x00 ||
		cc == 0x01 && Z == 0x01 ||
		cc == 0x02 && C == 0x00 ||
		cc == 0x03 && C == 0x01 {
		address := binary.LittleEndian.Uint16(ins[1:])
		gb.set16Reg(PC, address)
	} else {
		gb.regs[PC] += uint16(len(ins))
	}
}

func (gb *GameBoy) JR_e(ins [2]uint8) {
	jump := int8(ins[1])
	gb.regs[PC] += uint16(len(ins))
	gb.regs[PC] = uint16(int16(gb.regs[PC]) + int16(jump))
}

func (gb *GameBoy) JR_cc_e(ins [2]uint8) {
	Z := gb.getFlag(Z_FLAG)
	C := gb.getFlag(C_FLAG)
	cc := ins[0] >> 3 & 0x03
	if cc == 0x00 && Z == 0x00 ||
		cc == 0x01 && Z == 0x01 ||
		cc == 0x02 && C == 0x00 ||
		cc == 0x03 && C == 0x01 {
		jump := int8(ins[1])
		gb.regs[PC] = uint16(int16(gb.regs[PC]) + int16(jump))
	}
	gb.regs[PC] += uint16(len(ins))
}

func (gb *GameBoy) JP_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	gb.set16Reg(PC, address)
}

// checkAndSetZeroFlag sets Z_FLAG if out is zero otherwise it clears it.
func (gb *GameBoy) checkAndSetZeroFlag(out uint8) {
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
}

// _rotateWithC handles rotation with/without C and in multiple directions
func _rotateWithC(input uint8, c_input uint8, left bool, withC bool) (output uint8, c_output uint16) {

	// handle c_output
	c_output = CLEAR
	if left {
		// left
		if input&0x80 == 0x80 {
			c_output = SET
		}
	} else {
		// right
		if input&0x1 == 0x1 {
			c_output = SET
		}
	}

	// do rotation
	direction := 1
	if !left {
		direction = -1
	}
	output = bits.RotateLeft8(input, direction)

	// handle rotating in C
	if withC {
		if left {
			// left
			if c_input&0x1 == 0x1 {
				output = output | 0x1
			} else {
				output = output & 0xfe
			}
		} else {
			// right
			if c_input&0x1 == 0x1 {
				output = output | 0x80
			} else {
				output = output & 0x7f
			}
		}
	}

	return
}

// RLCA rotates the contents of register A to the left. The contents of bit 7
// are placed in both C and bit 0 of the result.
func (gb *GameBoy) RLCA(ins [1]uint8) {
	val := gb.get8Reg(A)
	out, c_out := _rotateWithC(val, 0, true, false)
	gb.set8Reg(A, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// RLA rotates the contents of register A to the left. The contents of bit 7
// are placed in C, and C is rotated into bit 0.
func (gb *GameBoy) RLA(ins [1]uint8) {
	val := gb.get8Reg(A)
	out, c_out := _rotateWithC(val, gb.getFlag(C_FLAG), true, true)
	gb.set8Reg(A, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// RRCA rotates the contents of register A to the right. The contents of bit 0
// are placed in both C and bit 7 of the result.
func (gb *GameBoy) RRCA(ins [1]uint8) {
	val := gb.get8Reg(A)
	out, c_out := _rotateWithC(val, 0, false, false)
	gb.set8Reg(A, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// RRA rotates the contents of register A to the right. The contents of bit 0
// are placed in C, and C is rotated into bit 7.
func (gb *GameBoy) RRA(ins [1]uint8) {
	val := gb.get8Reg(A)
	out, c_out := _rotateWithC(val, gb.getFlag(C_FLAG), false, true)
	gb.set8Reg(A, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// RLC_r rotates the contents of register r to the left.
func (gb *GameBoy) RLC_r(ins [1]uint8) {
	r := Reg8ID((ins[0] >> 3) & 0x07)
	val := gb.get8Reg(r)
	out, c_out := _rotateWithC(val, 0, true, false)
	gb.set8Reg(r, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// RLC_hl rotates the data stored at address (HL) to the left.
func (gb *GameBoy) RLC_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	val := gb.mainMemory.read(address)
	out, c_out := _rotateWithC(val, 0, true, false)
	gb.mainMemory.write(address, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// RL_r rotates the contents of register r to the left.
func (gb *GameBoy) RL_r(ins [1]uint8) {
	r := Reg8ID((ins[0] >> 3) & 0x07)
	val := gb.get8Reg(r)
	out, c_out := _rotateWithC(val, gb.getFlag(C_FLAG), true, true)
	gb.set8Reg(r, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// RL_hl rotates the data stored at address (HL) to the left.
func (gb *GameBoy) RL_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	val := gb.mainMemory.read(address)
	out, c_out := _rotateWithC(val, gb.getFlag(C_FLAG), true, true)
	gb.mainMemory.write(address, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// RRC_r rotates the contents of register r to the right.
func (gb *GameBoy) RRC_r(ins [1]uint8) {
	r := Reg8ID((ins[0] >> 3) & 0x07)
	val := gb.get8Reg(r)
	out, c_out := _rotateWithC(val, 0, false, false)
	gb.set8Reg(r, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// RRC_hl rotates the data stored at address (HL) to the right.
func (gb *GameBoy) RRC_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	val := gb.mainMemory.read(address)
	out, c_out := _rotateWithC(val, 0, false, false)
	gb.mainMemory.write(address, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// RR_r rotates the contents of register r to the right.
func (gb *GameBoy) RR_r(ins [1]uint8) {
	r := Reg8ID((ins[0] >> 3) & 0x07)
	val := gb.get8Reg(r)
	out, c_out := _rotateWithC(val, gb.getFlag(C_FLAG), false, true)
	gb.set8Reg(r, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// RR_hl rotates the data stored at address (HL) to the right.
func (gb *GameBoy) RR_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	val := gb.mainMemory.read(address)
	out, c_out := _rotateWithC(val, gb.getFlag(C_FLAG), false, true)
	gb.mainMemory.write(address, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// SLA_r shifts the contents of register r to the left
func (gb *GameBoy) SLA_r(ins [1]uint8) {
	r := Reg8ID((ins[0] >> 3) & 0x07)
	val := gb.get8Reg(r)

	// rotate and set bit 0 to 0 for shift
	out := bits.RotateLeft8(val, 1)
	out = out & 0xfe
	gb.set8Reg(r, out)

	// shift 0x80 into C
	if val&0x80 == 0x80 {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}

	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// SLA_hl shifts the data stored at address (HL) to the left
func (gb *GameBoy) SLA_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	val := gb.mainMemory.read(address)

	// rotate and set bit 0 to 0 for shift
	out := bits.RotateLeft8(val, 1)
	out = out & 0xfe
	gb.mainMemory.write(address, out)

	// shift 0x80 into C
	if val&0x80 == 0x80 {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}

	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// SRA_r shifts the contents of register r to the right
func (gb *GameBoy) SRA_r(ins [1]uint8) {
	r := Reg8ID((ins[0] >> 3) & 0x07)
	val := gb.get8Reg(r)

	// rotate and set bit 7 to its original value
	out := bits.RotateLeft8(val, -1)
	if val&0x80 == 0x80 {
		out = out | 0x80
	} else {
		out = out & 0x7f
	}
	gb.set8Reg(r, out)

	// shift 0x1 into C
	if val&0x1 == 0x1 {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}

	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// SRA_hl shifts the data stored at address (HL) to the right
func (gb *GameBoy) SRA_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	val := gb.mainMemory.read(address)

	// rotate and set bit 7 to its original value
	out := bits.RotateLeft8(val, -1)
	if val&0x80 == 0x80 {
		out = out | 0x80
	} else {
		out = out & 0x7f
	}
	gb.mainMemory.write(address, out)

	// shift 0x1 into C
	if val&0x1 == 0x1 {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}

	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// SRL r
func (gb *GameBoy) SRL_r(ins [1]uint8) {
	r := Reg8ID((ins[0] >> 3) & 0x07)
	val := gb.get8Reg(r)

	// rotate and set bit 7 to 0 for shift
	out := bits.RotateLeft8(val, -1)
	out = out & 0x7f
	gb.set8Reg(r, out)

	// shift 0x1 into C
	if val&0x1 == 0x1 {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}

	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// SRL_hl shifts the data stored at address (HL) to the right
func (gb *GameBoy) SRL_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	val := gb.mainMemory.read(address)

	// rotate and set bit 7 to 0 for shift
	out := bits.RotateLeft8(val, -1)
	out = out & 0x7f
	gb.mainMemory.write(address, out)

	// shift 0x1 into C
	if val&0x1 == 0x1 {
		gb.modifyFlag(C_FLAG, SET)
	} else {
		gb.modifyFlag(C_FLAG, CLEAR)
	}

	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// SWAP_r
func (gb *GameBoy) SWAP_r(ins [1]uint8) {
	r := Reg8ID((ins[0] >> 3) & 0x07)
	val := gb.get8Reg(r)
	// nibble swap
	ln := val & 0xf
	un := (val & 0xf0) >> 4
	out := (ln << 4) & un
	gb.set8Reg(r, out)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// SWAP_hl
func (gb *GameBoy) SWAP_hl(ins [1]uint8) {
	address := gb.get16Reg(HL)
	val := gb.mainMemory.read(address)
	// nibble swap
	ln := val & 0xf
	un := (val & 0xf0) >> 4
	out := (ln << 4) & un
	gb.mainMemory.write(address, out)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}
