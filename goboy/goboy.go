package main

import "image"

//Frances was here!

type GameBoy struct {
	rom              *GBROM // the ROM object
	mainMemory       *GBMem // GB main memory
	*Register               // register state
	interruptEnabled bool
	image            *image.RGBA // image to be displayed
}

type Reg8ID int
type Reg16ID int
type FlagId uint16

const (
	B Reg8ID = 0x00
	C Reg8ID = 0x01
	D Reg8ID = 0x02
	E Reg8ID = 0x03
	H Reg8ID = 0x04
	L Reg8ID = 0x05
	F Reg8ID = 0x06 // NOTE: it's not confirmed that 0x06 is actually the register code for F
	A Reg8ID = 0x07
)

const (
	BC Reg16ID = 0x00
	DE Reg16ID = 0x01
	HL Reg16ID = 0x02
	SP Reg16ID = 0x03
	PC Reg16ID = 0x04
	AF Reg16ID = 0x05
)

const (
	Z_FLAG FlagId = 0x0080
	N_FLAG FlagId = 0x0040
	H_FLAG FlagId = 0x0020
	C_FLAG FlagId = 0x0010
)

const (
	SET   = 0xffff
	CLEAR = 0x0000
)

type Register struct {
	/*
	   Register structure
	   +---------------+
	   |   B   |   C   |
	   +---------------+
	   |   D   |   E   |
	   +---------------+
	   |   H   |   L   |
	   +---------------+
	   |      SP       |
	   +---------------+
	   |      PC       |
	   +---------------+
	   |   A   |   F   |
	   +---------------+
	*/
	regs [6]uint16
}

func (r *Register) get16Reg(id Reg16ID) uint16 {
	return r.regs[id]
}

func (r *Register) set16Reg(id Reg16ID, value uint16) {
	r.regs[id] = value
}

func (r *Register) get8Reg(id Reg8ID) uint8 {
	if id == A { // Seperate logic because A and F are special snowflakes
		return uint8(r.regs[5] >> 8)
	} else if id == F {
		return uint8(r.regs[5])
	}
	block := id / 2 // block is the index of the 16bit version
	end := id % 2   // end indicates the high (0) or low (1) byte of the 16bit register
	if end == 0 {
		return uint8(r.regs[block] >> 8)
	} else {
		return uint8(r.regs[block])
	}
}

func (r *Register) set8Reg(id Reg8ID, value uint8) {
	if id == A { //Seperate logic because A and F are special snowflakes
		r.regs[5] &= 0x00ff
		r.regs[5] |= uint16(value) << 8
		return
	} else if id == F {
		r.regs[5] &= 0xff00
		r.regs[5] |= uint16(value)
		return
	}
	block := id / 2
	end := id % 2
	if end == 0 {
		r.regs[block] &= 0x00ff
		r.regs[block] |= uint16(value) << 8
	} else {
		r.regs[block] &= 0xff00
		r.regs[block] |= uint16(value)
	}
}

// modifyFlag sets the flag to a value which is 0 or 1
func (r *Register) modifyFlag(flag FlagId, value uint16) {
	r.regs[AF] &= ^uint16(flag)
	r.regs[AF] |= (uint16(flag) & value)
}

// getFlag returns the value of the flag
func (r *Register) getFlag(flag FlagId) uint8 {
	if r.regs[AF]&uint16(flag) > 0 {
		return 1
	} else {
		return 0
	}
}
