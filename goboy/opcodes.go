package main

import (
	"encoding/binary"
	"math/bits"
)

// Load r1 <- r2 1B
func (gb *GameBoy) LD_r_r(ins []uint8) {
	r1 := Reg8ID((ins[0] & 0x38) >> 3)
	r2 := Reg8ID(ins[0] & 0x07)
	gb.set8Reg(r1, gb.get8Reg(r2))
	gb.regs[PC] += uint16(len(ins))
}

// Load r <- n 2B
func (gb *GameBoy) LD_r_n(ins []uint8) {
	n := ins[1]
	r := Reg8ID(ins[0] >> 3)
	gb.set8Reg(r, n)
	gb.regs[PC] += uint16(len(ins))
}

// Load r <- (HL) 1B
func (gb *GameBoy) LD_r_hl(ins []uint8) {
	r := Reg8ID((ins[0] & 0x38) >> 3)
	address := gb.get16Reg(HL)
	gb.set8Reg(r, gb.mainMemory.read(address))
	gb.regs[PC] += uint16(len(ins))
}

// Load (HL) <- r 1B
func (gb *GameBoy) LD_hl_r(ins []uint8) {
	r := Reg8ID(ins[0] & 0x07)
	address := gb.get16Reg(HL)
	gb.mainMemory.write(address, gb.get8Reg(r))
	gb.regs[PC] += uint16(len(ins))
}

// Load (HL) <- n 2B
func (gb *GameBoy) LD_hl_n(ins []uint8) {
	n := ins[1]
	address := gb.get16Reg(HL)
	gb.mainMemory.write(address, n)
	gb.regs[PC] += uint16(len(ins))
}

// Load A <- (BC) 1B
func (gb *GameBoy) LD_a_bc(ins []uint8) {
	address := gb.get16Reg(BC)
	gb.set8Reg(A, gb.mainMemory.read(address))
	gb.regs[PC] += uint16(len(ins))
}

// Load A <- (DE) 1B
func (gb *GameBoy) LD_a_de(ins []uint8) {
	address := gb.get16Reg(DE)
	gb.set8Reg(A, gb.mainMemory.read(address))
	gb.regs[PC] += uint16(len(ins))
}

// LOAD A <- (0xff00 + C)
// TODO: Implement Unittest
func (gb *GameBoy) LD_a_c(ins []uint8) {
    address := uint16(gb.get8Reg(C)) + 0xff00
    value := gb.mainMemory.read(address)
    gb.set8Reg(A, value)
    gb.regs[PC] += uint16(len(ins))
}

// Load (0xff00 + C) <- A
// TODO: Implement Unittest
func (gb *GameBoy) LD_c_a(ins []uint8) {
    value := gb.get8Reg(A)
    address := uint16(gb.get8Reg(C)) + 0xff00
    gb.mainMemory.write(address, value)
    gb.regs[PC] += uint16(len(ins))
}

// Load A <- (0xff00 + n)
// TODO: Implement Unittest
func (gb *GameBoy) LD_a_n(ins []uint8) {
    address := uint16(ins[1]) + 0xff00
    value := gb.mainMemory.read(address)
    gb.set8Reg(A, value)
    gb.regs[PC] += uint16(len(ins))
}

// Load (0xff00 + n) <- A
// TODO: Implement Unittest
func (gb *GameBoy) LD_n_a(ins []uint8) {
    address := uint16(ins[1]) + 0xff00
    value := gb.get8Reg(A)
    gb.mainMemory.write(address, value)
    gb.regs[PC] += uint16(len(ins))
}

// Load A <- (nn) 3B
func (gb *GameBoy) LD_a_nn(ins []uint8) {
	address := binary.LittleEndian.Uint16(ins[1:])
	gb.set8Reg(A, gb.mainMemory.read(address))
	gb.regs[PC] += uint16(len(ins))
}

// Load (BC) <- A 1B
func (gb *GameBoy) LD_bc_a(ins []uint8) {
	address := gb.get16Reg(BC)
	gb.mainMemory.write(address, gb.get8Reg(A))
	gb.regs[PC] += uint16(len(ins))
}

// Load (DE) <- A 1B
func (gb *GameBoy) LD_de_a(ins []uint8) {
	address := gb.get16Reg(DE)
	gb.mainMemory.write(address, gb.get8Reg(A))
	gb.regs[PC] += uint16(len(ins))
}

// Load (nn) <- A 3B
func (gb *GameBoy) LD_nn_a(ins []uint8) {
	address := binary.LittleEndian.Uint16(ins[1:])
	gb.mainMemory.write(address, gb.get8Reg(A))
	gb.regs[PC] += uint16(len(ins))
}

// Load A <- (HL); HL++
// TODO: Implement Unittest
func (gb *GameBoy) LD_a_hli(ins []uint8) {
    address := gb.get16Reg(HL)
    value := gb.mainMemory.read(address)
    gb.set8Reg(A, value)
    gb.set16Reg(HL, address + 0x1)
    gb.regs[PC] += uint16(len(ins))
}

// Load A <- (HL); HL--
// TODO: Implement Unittest
func (gb *GameBoy) LD_a_hld(ins []uint8) {
    address := gb.get16Reg(HL)
    value := gb.mainMemory.read(address)
    gb.set8Reg(A, value)
    gb.set16Reg(HL, address - 0x1)
    gb.regs[PC] += uint16(len(ins))
}

// Load (HL) <- A; HL++
// TODO: Implement Unittest
func (gb *GameBoy) LD_hli_a(ins []uint8) {
    value := gb.get8Reg(A)
    address := gb.get16Reg(HL)
    gb.mainMemory.write(address, value)
    gb.set16Reg(HL, address + 0x1)
    gb.regs[PC] += uint16(len(ins))
}

// Load (HL) <- A; HL--
// TODO: Implement Unittest
func (gb *GameBoy) LD_hld_a(ins []uint8) {
    value := gb.get8Reg(A)
    address := gb.get16Reg(HL)
    gb.mainMemory.write(address, value)
    gb.set16Reg(HL, address - 0x1)
    gb.regs[PC] += uint16(len(ins))
}

/* 16 BIT LOADS */

// LD dd <- nn 3B
func (gb *GameBoy) LD_dd_nn(ins []uint8) {
	reg := Reg16ID(ins[0] >> 4)
	nn := binary.LittleEndian.Uint16(ins[1:])
	gb.set16Reg(reg, nn)
	gb.regs[PC] += uint16(len(ins))
}

// LD SP <- HL 1B
func (gb *GameBoy) LD_sp_hl(ins []uint8) {
	gb.set16Reg(SP, gb.get16Reg(HL))
	gb.regs[PC] += uint16(len(ins))
}

// PUSH qq 1B
func (gb *GameBoy) PUSH_qq(ins []uint8) {
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

// POP qq 1B
func (gb *GameBoy) POP_qq(ins []uint8) {
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

// Load HL <- SP + e
// TODO: Implement Unittest
func (gb *GameBoy) LDHL_sp_e(ins []uint8) {
    aVal := gb.get16Reg(SP)
    bVal := uint16(int8(ins[1]))
    out := aVal + bVal
    gb.set16Reg(HL, out)
    gb.modifyFlag(Z_FLAG, CLEAR)
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

// Load (nn) <- SP_lower; (nn + 1) <- SP_upper
// TODO: Implement Unittest
func (gb *GameBoy) LD_nn_sp(ins []uint8) {
    sp := gb.get16Reg(SP)
    sp_lb := uint8(sp)
    sp_hb := uint8(sp >> 8)
    address := binary.LittleEndian.Uint16(ins[1:])
    gb.mainMemory.write(address, sp_lb)
    gb.mainMemory.write(address + 0x1, sp_hb)
    gb.regs[PC] += uint16(len(ins))
}

// ADD A, r 1B
func (gb *GameBoy) ADD_a_r(ins []uint8) {
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

// ADD A, n 2B
func (gb *GameBoy) ADD_a_n(ins []uint8) {
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

// ADD A, (HL) 1B
func (gb *GameBoy) ADD_a_hl(ins []uint8) {
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

// ADC A, r 1B
func (gb *GameBoy) ADC_a_r(ins []uint8) {
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

// ADC A, n 2B
func (gb *GameBoy) ADC_a_n(ins []uint8) {
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

// ADC A, (HL) 1B
func (gb *GameBoy) ADC_a_hl(ins []uint8) {
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

// SUB A, r 1B
func (gb *GameBoy) SUB_a_r(ins []uint8) {
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

// SUB A, n 2B
func (gb *GameBoy) SUB_a_n(ins []uint8) {
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

// SUB A, (HL) 1B
func (gb *GameBoy) SUB_a_hl(ins []uint8) {
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

// SBC A, r 1B
func (gb *GameBoy) SBC_a_r(ins []uint8) {
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

// SBC A, n 2B
func (gb *GameBoy) SBC_a_n(ins []uint8) {
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

// SBC A, (HL) 1B
func (gb *GameBoy) SBC_a_hl(ins []uint8) {
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

// AND A, r 1B
func (gb *GameBoy) AND_a_r(ins []uint8) {
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

// AND A, n 2B
func (gb *GameBoy) AND_a_n(ins []uint8) {
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

// AND A, (HL) 1B
func (gb *GameBoy) AND_a_hl(ins []uint8) {
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

// OR A, r 1B
func (gb *GameBoy) OR_a_r(ins []uint8) {
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

// OR A, n 2B
func (gb *GameBoy) OR_a_n(ins []uint8) {
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

// OR A, (HL) 1B
func (gb *GameBoy) OR_a_hl(ins []uint8) {
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

// XOR A, r 1B
func (gb *GameBoy) XOR_a_r(ins []uint8) {
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

// XOR A, n 2B
func (gb *GameBoy) XOR_a_n(ins []uint8) {
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

// XOR A, (HL) 1B
func (gb *GameBoy) XOR_a_hl(ins []uint8) {
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

// CP A, r 1B
func (gb *GameBoy) CP_a_r(ins []uint8) {
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

// CP A, n 2B
func (gb *GameBoy) CP_a_n(ins []uint8) {
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

// CP A,(HL) 1B
func (gb *GameBoy) CP_a_hl(ins []uint8) {
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

// INC r 1B
func (gb *GameBoy) INC_r(ins []uint8) {
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

// INC (HL) 1B
func (gb *GameBoy) INC_hl(ins []uint8) {
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

// DEC r 1B
func (gb *GameBoy) DEC_r(ins []uint8) {
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

// DEC (HL) 1B
func (gb *GameBoy) DEC_hl(ins []uint8) {
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

// ADD HL ss 1B
func (gb *GameBoy) ADD_hl_ss(ins []uint8) {
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

// ADD SP e 2B
func (gb *GameBoy) ADD_sp_e(ins []uint8) {
	aVal := gb.get16Reg(SP)
	bVal := uint16(int8(ins[1]))
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

// 16 bit INC 1B

/*Retrieve value of ss
    Shift instruction array to the right by 4
    Clear non-ss bits

Get value of ss register

Add 1 to register value

Store incremented value back into register ss*/

func (gb *GameBoy) INC_ss(ins []uint8) {
	ss := Reg16ID((ins[0] >> 4) & 0x03)
	ssVal := gb.get16Reg(ss)
	out := ssVal + 0x01
	gb.set16Reg(ss, out)
	gb.regs[PC] += uint16(len(ins))
}

// 16 bit DEC 1B
func (gb *GameBoy) DEC_ss(ins []uint8) {
	ss := Reg16ID((ins[0] >> 4) & 0x03)
	ssVal := gb.get16Reg(ss)
	out := ssVal - 0x01
	gb.set16Reg(ss, out)
	gb.regs[PC] += uint16(len(ins))
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
// 1B
func (gb *GameBoy) RLCA(ins []uint8) {
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
// 1B
func (gb *GameBoy) RLA(ins []uint8) {
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
// 1B
func (gb *GameBoy) RRCA(ins []uint8) {
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
// 1B
func (gb *GameBoy) RRA(ins []uint8) {
	val := gb.get8Reg(A)
	out, c_out := _rotateWithC(val, gb.getFlag(C_FLAG), false, true)
	gb.set8Reg(A, out)
	gb.modifyFlag(C_FLAG, uint16(c_out))
	gb.checkAndSetZeroFlag(out)
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.modifyFlag(N_FLAG, CLEAR)
	gb.regs[PC] += uint16(len(ins))
}

// RLC_r rotates the contents of register r to the left. 1B
func (gb *GameBoy) RLC_r(ins []uint8) {
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

// RLC_hl rotates the data stored at address (HL) to the left. 1B
func (gb *GameBoy) RLC_hl(ins []uint8) {
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

// RL_r rotates the contents of register r to the left. 1B
func (gb *GameBoy) RL_r(ins []uint8) {
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

// RL_hl rotates the data stored at address (HL) to the left. 1B
func (gb *GameBoy) RL_hl(ins []uint8) {
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

// RRC_r rotates the contents of register r to the right. 1B
func (gb *GameBoy) RRC_r(ins []uint8) {
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

// RRC_hl rotates the data stored at address (HL) to the right. 1B
func (gb *GameBoy) RRC_hl(ins []uint8) {
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

// RR_r rotates the contents of register r to the right. 1B
func (gb *GameBoy) RR_r(ins []uint8) {
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

// RR_hl rotates the data stored at address (HL) to the right. 1B
func (gb *GameBoy) RR_hl(ins []uint8) {
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

// SLA_r shifts the contents of register r to the left 1B
func (gb *GameBoy) SLA_r(ins []uint8) {
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

// SLA_hl shifts the data stored at address (HL) to the left 1B
func (gb *GameBoy) SLA_hl(ins []uint8) {
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

// SRA_r shifts the contents of register r to the right 1B
func (gb *GameBoy) SRA_r(ins []uint8) {
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

// SRA_hl shifts the data stored at address (HL) to the right 1B
func (gb *GameBoy) SRA_hl(ins []uint8) {
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

// SRL r 1B
func (gb *GameBoy) SRL_r(ins []uint8) {
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

// SRL_hl shifts the data stored at address (HL) to the right 1B
func (gb *GameBoy) SRL_hl(ins []uint8) {
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

// SWAP_r 1B
func (gb *GameBoy) SWAP_r(ins []uint8) {
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

// SWAP_hl 1B
func (gb *GameBoy) SWAP_hl(ins []uint8) {
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

// BIT b r 2B
func (gb *GameBoy) BIT_b_r(ins []uint8) {
	r := Reg8ID(ins[1] & 0x07)
	b := (ins[1] >> 3) & 0x07
	value := ((^gb.get8Reg(r)) >> b) & 0x01 // extract the bth bit from ~r
	gb.modifyFlag(H_FLAG, SET)
	gb.modifyFlag(N_FLAG, CLEAR)
	if value == 1 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.regs[PC] += uint16(len(ins))
}

// BIT b (HL) 2B
func (gb *GameBoy) BIT_b_hl(ins []uint8) {
	b := (ins[1] >> 3) & 0x07
	address := gb.get16Reg(HL)
	value := ((^gb.mainMemory.read(address)) >> b) & 0x01 // extract the bth bit from ~(HL)
	gb.modifyFlag(H_FLAG, SET)
	gb.modifyFlag(N_FLAG, CLEAR)
	if value == 1 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.regs[PC] += uint16(len(ins))
}

// SET b r 2B
func (gb *GameBoy) SET_b_r(ins []uint8) {
	r := Reg8ID(ins[1] & 0x07)
	b := (ins[1] >> 3) & 0x07
	value := gb.get8Reg(r) | (0x01 << b)
	gb.set8Reg(r, value)
	gb.regs[PC] += uint16(len(ins))
}

// SET b (HL) 2B
func (gb *GameBoy) SET_b_hl(ins []uint8) {
	b := (ins[1] >> 3) & 0x07
	address := gb.get16Reg(HL)
	value := gb.mainMemory.read(address) | (0x01 << b)
	gb.mainMemory.write(address, value)
	gb.regs[PC] += uint16(len(ins))
}

// RES b r 2B
func (gb *GameBoy) RES_b_r(ins []uint8) {
	r := Reg8ID(ins[1] & 0x07)
	b := (ins[1] >> 3) & 0x07
	value := gb.get8Reg(r) & ^(0x01 << b)
	gb.set8Reg(r, value)
	gb.regs[PC] += uint16(len(ins))
}

// RES b (HL) 2B
func (gb *GameBoy) RES_b_hl(ins []uint8) {
	b := (ins[1] >> 3) & 0x07
	address := gb.get16Reg(HL)
	value := gb.mainMemory.read(address) & ^(0x01 << b)
	gb.mainMemory.write(address, value)
	gb.regs[PC] += uint16(len(ins))
}

// JP nn 3B
func (gb *GameBoy) JP_nn(ins []uint8) {
	address := binary.LittleEndian.Uint16(ins[1:])
	gb.set16Reg(PC, address)
}

// JP cc nn 3B
func (gb *GameBoy) JP_cc_nn(ins []uint8) {
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

// JR e 2B
func (gb *GameBoy) JR_e(ins []uint8) {
	jump := int8(ins[1])
	gb.regs[PC] += uint16(len(ins))
	gb.regs[PC] = uint16(int16(gb.regs[PC]) + int16(jump))
}

// JR cc e 2B
func (gb *GameBoy) JR_cc_e(ins []uint8) {
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

// JP (HL) 1B
func (gb *GameBoy) JP_hl(ins []uint8) {
	address := gb.get16Reg(HL)
	gb.set16Reg(PC, address)
}

// CALL nn 3B
func (gb *GameBoy) CALL_nn(ins []uint8) {
	ret_pc := gb.regs[PC] + uint16(len(ins))
	gb.set16Reg(PC, binary.LittleEndian.Uint16(ins[1:]))
	gb.regs[SP]--
	gb.mainMemory.write(gb.regs[SP], uint8(ret_pc>>8))
	gb.regs[SP]--
	gb.mainMemory.write(gb.regs[SP], uint8(ret_pc&0x00ff))
}

// CALL cc nn 3B
func (gb *GameBoy) CALL_cc_nn(ins []uint8) {
	Z := gb.getFlag(Z_FLAG)
	C := gb.getFlag(C_FLAG)
	cc := ins[0] >> 3 & 0x03
	if cc == 0x00 && Z == 0x00 ||
		cc == 0x01 && Z == 0x01 ||
		cc == 0x02 && C == 0x00 ||
		cc == 0x03 && C == 0x01 {
		ret_pc := gb.regs[PC] + uint16(len(ins))
		gb.set16Reg(PC, binary.LittleEndian.Uint16(ins[1:]))
		gb.regs[SP]--
		gb.mainMemory.write(gb.regs[SP], uint8(ret_pc>>8))
		gb.regs[SP]--
		gb.mainMemory.write(gb.regs[SP], uint8(ret_pc&0x00ff))
	} else {
		gb.regs[PC] += uint16(len(ins))
	}
}

// RET 1B
func (gb *GameBoy) RET(ins []uint8) {
	address_lsb := gb.mainMemory.read(gb.get16Reg(SP))
	gb.regs[SP]++
	address_msb := gb.mainMemory.read(gb.get16Reg(SP))
	gb.regs[SP]++
	gb.set16Reg(PC, binary.LittleEndian.Uint16([]uint8{address_lsb, address_msb}))
}

// RETI 1B
func (gb *GameBoy) RETI(ins []uint8) {
	gb.interruptEnabled = true
	address_lsb := gb.mainMemory.read(gb.get16Reg(SP))
	gb.regs[SP]++
	address_msb := gb.mainMemory.read(gb.get16Reg(SP))
	gb.regs[SP]++
	gb.set16Reg(PC, binary.LittleEndian.Uint16([]uint8{address_lsb, address_msb}))
}

// RET cc 1B
func (gb *GameBoy) RET_cc(ins []uint8) {
	Z := gb.getFlag(Z_FLAG)
	C := gb.getFlag(C_FLAG)
	cc := ins[0] >> 3 & 0x03
	if cc == 0x00 && Z == 0x00 ||
		cc == 0x01 && Z == 0x01 ||
		cc == 0x02 && C == 0x00 ||
		cc == 0x03 && C == 0x01 {
		address_lsb := gb.mainMemory.read(gb.get16Reg(SP))
		gb.regs[SP]++
		address_msb := gb.mainMemory.read(gb.get16Reg(SP))
		gb.regs[SP]++
		gb.set16Reg(PC, binary.LittleEndian.Uint16([]uint8{address_lsb, address_msb}))
	} else {
		gb.regs[PC] += uint16(len(ins))
	}
}

// RST 1B
func (gb *GameBoy) RST(ins []uint8) {
	t := uint16((ins[0] >> 3) & 0x07)
	ret_pc := gb.regs[PC] + uint16(len(ins))
	gb.regs[SP]--
	gb.mainMemory.write(gb.regs[SP], uint8(ret_pc>>8))
	gb.regs[SP]--
	gb.mainMemory.write(gb.regs[SP], uint8(ret_pc&0x00ff))
	gb.set16Reg(PC, 0x0008*t)
}

// DAA 1B
func (gb *GameBoy) DAA(ins []uint8) {
	num := gb.get8Reg(A)
	lsn := num & 0x0f        // least significant nibble
	msn := (num & 0xf0) >> 4 // most significant nibble
	correction := uint8(0x00)
	if gb.getFlag(N_FLAG) == 0 { // Addition case
		if gb.getFlag(H_FLAG) == 1 || lsn > 0x09 {
			correction += 0x06
		}
		if gb.getFlag(C_FLAG) == 1 || msn > 0x09 ||
			(msn == 0x09 && lsn > 0x09) {
			correction += 0x60
			gb.modifyFlag(C_FLAG, SET)
		}
	} else { // Subtraction case
		if gb.getFlag(H_FLAG) == 1 {
			correction -= 0x06
		}
		if gb.getFlag(C_FLAG) == 1 {
			correction -= 0x60
			gb.modifyFlag(C_FLAG, SET)
		}
	}
	out := num + correction
	if out == 0 {
		gb.modifyFlag(Z_FLAG, SET)
	} else {
		gb.modifyFlag(Z_FLAG, CLEAR)
	}
	gb.modifyFlag(H_FLAG, CLEAR)
	gb.set8Reg(A, out)
	gb.regs[PC] += uint16(len(ins))
}

// CPL 1B
func (gb *GameBoy) CPL(ins []uint8) {
	gb.set8Reg(A, ^gb.get8Reg(A))
	gb.regs[PC] += uint16(len(ins))
}

// NOP 1B
func (gb *GameBoy) NOP(ins []uint8) {
	gb.regs[PC] += uint16(len(ins))
}
