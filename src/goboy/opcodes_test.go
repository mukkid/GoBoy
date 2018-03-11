package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func initGameboy() *GameBoy {
	return &GameBoy{
		Register:   &Register{},
		mainMemory: &GBMem{},
	}
}

func TestLD_r_r(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(B, 0xfe)
	gb.LD_r_r([1]uint8{0x50}) // LD D, B
	assert.Equal(t, gb.get8Reg(B), uint8(0xfe))
	assert.Equal(t, gb.get8Reg(D), uint8(0xfe))
	assert.Equal(t, gb.get16Reg(PC), uint16(0x0001))
}

func TestLD_r_n(t *testing.T) {
	gb := initGameboy()
	gb.LD_r_n([2]uint8{0x0e, 0xad}) // LD C, 0xad
	assert.Equal(t, gb.get8Reg(C), uint8(0xad))
	assert.Equal(t, gb.get16Reg(PC), uint16(0x0002))
}

func TestLD_r_hl(t *testing.T) {
	gb := initGameboy()
	gbROM := newGBROM()
	gbROM.rom[0x1234] = 0x42
	gb.mainMemory.cartridge = gbROM
	memValue := gb.mainMemory.read(0x1234)
	assert.Equal(t, memValue, uint8(0x42))
	gb.set16Reg(HL, 0x1234)
	gb.LD_r_hl([1]uint8{0x5e}) // LD E (HL)
	assert.Equal(t, gb.get8Reg(E), uint8(0x42))
}

func TestLD_hl_r(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(B, 0x42)
	gb.set16Reg(HL, 0xff85)
	gb.LD_hl_r([1]uint8{0x70}) // LD (HL) B
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0x42))
}

func TestLD_hl_n(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(HL, 0xff85)
	gb.LD_hl_n([2]uint8{0x36, 0x42}) // LD (HL) 0x42
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0x42))
}

func TestLD_a_bc(t *testing.T) {
	gb := initGameboy()
	gbROM := newGBROM()
	gbROM.rom[0x1234] = 0x42
	gb.mainMemory.cartridge = gbROM

	memValue := gb.mainMemory.read(0x1234)
	assert.Equal(t, memValue, uint8(0x42))
	gb.set16Reg(BC, 0x1234)
	gb.LD_a_bc([1]uint8{0x0a}) // LD A (BC)
	assert.Equal(t, gb.get8Reg(A), uint8(0x42))
}

func TestLD_a_de(t *testing.T) {
	gb := initGameboy()
	gbROM := newGBROM()
	gbROM.rom[0x1234] = 0x42
	gb.mainMemory.cartridge = gbROM
	memValue := gb.mainMemory.read(0x1234)
	assert.Equal(t, memValue, uint8(0x42))
	gb.set16Reg(DE, 0x1234)
	gb.LD_a_de([1]uint8{0x1a}) // LD A (DE)
	assert.Equal(t, gb.get8Reg(A), uint8(0x42))
}

func TestLD_a_nn(t *testing.T) {
	gb := initGameboy()
	gbROM := newGBROM()
	gbROM.rom[0x1234] = 0x42
	gb.mainMemory.cartridge = gbROM

	memValue := gb.mainMemory.read(0x1234)
	assert.Equal(t, memValue, uint8(0x42))
	gb.LD_a_nn([3]uint8{0x3a, 0x34, 0x12}) // LD A (0x1234)
	assert.Equal(t, gb.get8Reg(A), uint8(0x42))
	assert.Equal(t, gb.get16Reg(PC), uint16(0x3))
}

func TestLD_bc_a(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0x42)
	gb.set16Reg(BC, 0xff85)
	gb.LD_bc_a([1]uint8{0x02}) // LD 0xff85 A
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0x42))
}

func TestLD_de_a(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0x42)
	gb.set16Reg(DE, 0xff85)
	gb.LD_de_a([1]uint8{0x12}) // LD 0xff85 A
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0x42))
}

func TestLD_nn_a(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0x42)
	gb.LD_nn_a([3]uint8{0x32, 0x85, 0xff}) // LD 0xff85 A
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0x42))
}

/* 16 BIT LOAD TESTS */

func TestLD_dd_nn(t *testing.T) {
	gb := initGameboy()
	gb.LD_dd_nn([3]uint8{0x21, 0xcd, 0xab}) // LD HL 0xabcd
	assert.Equal(t, gb.get16Reg(HL), uint16(0xabcd))
}

func TestLD_hl_nn(t *testing.T) {
	gb := initGameboy()

	gbROM := newGBROM()
	gbROM.rom[0x0123] = 0x37
	gbROM.rom[0x0124] = 0xa1
	gb.mainMemory.cartridge = gbROM
	gb.LD_hl_nn([3]uint8{0x2a, 0x23, 0x01}) // LD HL <- (0x0123)
	assert.Equal(t, gb.get16Reg(HL), uint16(0xa137))
}

func TestLD_dd_NN(t *testing.T) {
	gb := initGameboy()
	gbROM := newGBROM()
	gbROM.rom[0x0123] = 0xcd
	gbROM.rom[0x0124] = 0xab
	gb.mainMemory.cartridge = gbROM
	gb.LD_dd_NN([4]uint8{0xed, 0x5b, 0x23, 0x01}) // LD DE (0x0123)
	assert.Equal(t, gb.get16Reg(DE), uint16(0xabcd))
}

func TestLD_nn_hl(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(HL, 0x1234)
	gb.LD_nn_hl([3]uint8{0x22, 0x85, 0xff}) // LD (0xff85) HL
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0x34))
	assert.Equal(t, gb.mainMemory.read(0xff86), uint8(0x12))
}

func TestLD_nn_dd(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(DE, 0x1234)
	gb.LD_nn_dd([4]uint8{0xed, 0x53, 0x85, 0xff}) // LD 0xff85 DE
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0x34))
	assert.Equal(t, gb.mainMemory.read(0xff86), uint8(0x12))
}

func TestLD_sp_hl(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(HL, 0x1234)
	gb.LD_sp_hl([1]uint8{0xf9}) // LD SP HL (HL = 0x1234)
	assert.Equal(t, gb.get16Reg(SP), uint16(0x1234))
}

func TestPUSH_qq(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(HL, 0x1234)
	gb.set16Reg(SP, 0xff85)
	gb.PUSH_qq([1]uint8{0xe5}) // PUSH HL
	assert.Equal(t, gb.mainMemory.read(0xff84), uint8(0x12))
	assert.Equal(t, gb.mainMemory.read(0xff83), uint8(0x34))
}

func TestPOP_qq(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(SP, 0xff83)
	gb.mainMemory.write(0xff83, 0x34)
	gb.mainMemory.write(0xff84, 0x12)
	gb.POP_qq([1]uint8{0xf1}) // POP AF
	assert.Equal(t, gb.get16Reg(AF), uint16(0x1234))
}

/* ALU TESTS */

func TestADD_a_r(t *testing.T) {
	gb := initGameboy()
	// N_FLAG is always reset
	gb.ADD_a_r([1]uint8{0x80}) // ADD A B
	assert.Equal(t, gb.get8Reg(A), uint8(0x0))
	assert.Equal(t, gb.get8Reg(F), uint8(0x80)) // Z_FLAG is set

	gb.set8Reg(A, 0x0a)
	gb.set8Reg(B, 0x0c)
	gb.ADD_a_r([1]uint8{0x80}) // ADD A B
	assert.Equal(t, gb.get8Reg(A), uint8(0x16))
	assert.Equal(t, gb.get8Reg(F), uint8(0x20)) // H_FLAG is set

	gb.set8Reg(A, 0xf0)
	gb.set8Reg(B, 0xf0)
	gb.ADD_a_r([1]uint8{0x80}) // ADD A B
	assert.Equal(t, gb.get8Reg(A), uint8(0xe0))
	assert.Equal(t, gb.get8Reg(F), uint8(0x10)) // C_FLAG is set
}

func TestADD_a_n(t *testing.T) {
	// N_FLAG is always reset
	gb := initGameboy()
	gb.ADD_a_n([2]uint8{0xc6, 0x00}) // ADD A 0x00
	assert.Equal(t, gb.get8Reg(A), uint8(0x0))
	assert.Equal(t, gb.get8Reg(F), uint8(0x80)) // Z_FLAG is set

	gb.set8Reg(A, 0x0a)
	gb.ADD_a_n([2]uint8{0xc6, 0x0c}) // ADD A 0x0c
	assert.Equal(t, gb.get8Reg(A), uint8(0x16))
	assert.Equal(t, gb.get8Reg(F), uint8(0x20)) // H_FLAG is set

	gb.set8Reg(A, 0xf0)
	gb.ADD_a_n([2]uint8{0xc6, 0xf0}) // ADD A 0xf0
	assert.Equal(t, gb.get8Reg(A), uint8(0xe0))
	assert.Equal(t, gb.get8Reg(F), uint8(0x10)) // C_FLAG is set
}

func TestADD_a_hl(t *testing.T) {
	// N_FLAG is always reset
	gb := initGameboy()
	gbROM := newGBROM()
	gb.mainMemory.cartridge = gbROM
	gb.ADD_a_hl([1]uint8{0x86}) // ADD A HL
	assert.Equal(t, gb.get8Reg(A), uint8(0x0))
	assert.Equal(t, gb.get8Reg(F), uint8(0x80)) // Z_FLAG is set
	gb.set8Reg(A, 0x0a)
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0x0c)
	gb.ADD_a_hl([1]uint8{0x86}) // ADD A HL
	assert.Equal(t, gb.get8Reg(A), uint8(0x16))
	assert.Equal(t, gb.get8Reg(F), uint8(0x20)) // H_FLAG is set

	gb.set8Reg(A, 0xf0)
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0xf0)
	gb.ADD_a_hl([1]uint8{0x86}) // ADD A HL
	assert.Equal(t, gb.get8Reg(A), uint8(0xe0))
	assert.Equal(t, gb.get8Reg(F), uint8(0x10)) // C_FLAG is set
}

func TestADC_a_r(t *testing.T) {
	gb := initGameboy()
	// N_FLAG is always reset
	gb.ADC_a_r([1]uint8{0x88}) // ADC A B
	assert.Equal(t, gb.get8Reg(A), uint8(0x0))
	assert.Equal(t, gb.get8Reg(F), uint8(0x80)) // Z_FLAG is set

	gb.set8Reg(A, 0x0a)
	gb.set8Reg(B, 0x0c)
	gb.set8Reg(F, 0x10)
	gb.ADC_a_r([1]uint8{0x88}) // ADC A B
	assert.Equal(t, gb.get8Reg(A), uint8(0x17))
	assert.Equal(t, gb.get8Reg(F), uint8(0x20)) // H_FLAG is set

	gb.set8Reg(A, 0xf0)
	gb.set8Reg(B, 0xf0)
	gb.set8Reg(F, 0x10)
	gb.ADC_a_r([1]uint8{0x88}) // ADC A B
	assert.Equal(t, gb.get8Reg(A), uint8(0xe1))
	assert.Equal(t, gb.get8Reg(F), uint8(0x10)) // C_FLAG is set
}

func TestADC_a_n(t *testing.T) {
	// N_FLAG is always reset
	gb := initGameboy()
	gb.ADC_a_n([2]uint8{0xce, 0x00}) // ADC A 0x00
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0x80)) // Z_FLAG is set

	gb.set8Reg(A, 0x0a)
	gb.set8Reg(F, 0x10)
	gb.ADC_a_n([2]uint8{0xce, 0x0c}) // ADC A 0x0c
	assert.Equal(t, gb.get8Reg(A), uint8(0x17))
	assert.Equal(t, gb.get8Reg(F), uint8(0x20)) // H_FLAG is set

	gb.set8Reg(A, 0xf0)
	gb.set8Reg(F, 0x10)
	gb.ADC_a_n([2]uint8{0xce, 0xf0}) // ADC A 0xf0
	assert.Equal(t, gb.get8Reg(A), uint8(0xe1))
	assert.Equal(t, gb.get8Reg(F), uint8(0x10)) // C_FLAG is set
}

func TestADC_a_hl(t *testing.T) {
	// N_FLAG is always reset
	gb := initGameboy()
	gbROM := newGBROM()
	gb.mainMemory.cartridge = gbROM
	gb.ADC_a_hl([1]uint8{0x8e}) // ADC A HL
	assert.Equal(t, gb.get8Reg(A), uint8(0x0))
	assert.Equal(t, gb.get8Reg(F), uint8(0x80)) // Z_FLAG is set

	gb.set8Reg(A, 0x0a)
	gb.set8Reg(F, 0x10)
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0x0c)
	gb.ADC_a_hl([1]uint8{0x8e}) // ADC A HL
	assert.Equal(t, gb.get8Reg(A), uint8(0x17))
	assert.Equal(t, gb.get8Reg(F), uint8(0x20)) // H_FLAG is set

	gb.set8Reg(A, 0xf0)
	gb.set8Reg(F, 0x10)
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0xf0)
	gb.ADC_a_hl([1]uint8{0x8e}) // ADC A HL
	assert.Equal(t, gb.get8Reg(A), uint8(0xe1))
	assert.Equal(t, gb.get8Reg(F), uint8(0x10)) // C_FLAG is set
}

func TestSUB_a_r(t *testing.T) {
	gb := initGameboy()
	gb.SUB_a_r([1]uint8{0x90})
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0xc0)) // Z_FLAG and N_FLAG is set

	gb.set8Reg(A, 0x02)
	gb.set8Reg(B, 0x01)
	gb.SUB_a_r([1]uint8{0x90})
	assert.Equal(t, gb.get8Reg(A), uint8(0x01))
	assert.Equal(t, gb.get8Reg(F), uint8(0x40)) // N_FLAG is set

	gb.set8Reg(A, 0x87)
	gb.set8Reg(B, 0x0f)
	gb.SUB_a_r([1]uint8{0x90})
	assert.Equal(t, gb.get8Reg(A), uint8(0x78))
	assert.Equal(t, gb.get8Reg(F), uint8(0x60)) // N_FLAG and H_FLAG is set

	gb.set8Reg(A, 0x0f)
	gb.set8Reg(B, 0xf0)
	gb.SUB_a_r([1]uint8{0x90})
	assert.Equal(t, gb.get8Reg(A), uint8(0x1f))
	assert.Equal(t, gb.get8Reg(F), uint8(0x50)) // N_FLAG and C_FLAG is set
}

func TestSUB_a_n(t *testing.T) {
	gb := initGameboy()
	gb.SUB_a_n([2]uint8{0xd6, 0x00})
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0xc0)) // Z_FLAG and N_FLAG is set

	gb.set8Reg(A, 0x02)
	gb.SUB_a_n([2]uint8{0xd6, 0x01})
	assert.Equal(t, gb.get8Reg(A), uint8(0x01))
	assert.Equal(t, gb.get8Reg(F), uint8(0x40)) // N_FLAG is set

	gb.set8Reg(A, 0x87)
	gb.SUB_a_n([2]uint8{0xd6, 0x0f})
	assert.Equal(t, gb.get8Reg(A), uint8(0x78))
	assert.Equal(t, gb.get8Reg(F), uint8(0x60)) // N_FLAG and H_FLAG is set

	gb.set8Reg(A, 0x0f)
	gb.SUB_a_n([2]uint8{0xd6, 0xf0})
	assert.Equal(t, gb.get8Reg(A), uint8(0x1f))
	assert.Equal(t, gb.get8Reg(F), uint8(0x50)) // N_FLAG and C_FLAG is set
}

func TestSUB_a_hl(t *testing.T) {
	gb := initGameboy()
	gbROM := newGBROM()
	gb.mainMemory.cartridge = gbROM
	gb.SUB_a_hl([1]uint8{0x96})
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0xc0)) // Z_FLAG and N_FLAG is set

	gb.set8Reg(A, 0x02)
	gb.mainMemory.write(0xff85, 0x01)
	gb.set16Reg(HL, 0xff85)
	gb.SUB_a_hl([1]uint8{0x96})
	assert.Equal(t, gb.get8Reg(A), uint8(0x01))
	assert.Equal(t, gb.get8Reg(F), uint8(0x40)) // N_FLAG is set

	gb.set8Reg(A, 0x87)
	gb.mainMemory.write(0xff85, 0x0f)
	gb.set16Reg(HL, 0xff85)
	gb.SUB_a_hl([1]uint8{0x96})
	assert.Equal(t, gb.get8Reg(A), uint8(0x78))
	assert.Equal(t, gb.get8Reg(F), uint8(0x60)) // N_FLAG and H_FLAG is set

	gb.set8Reg(A, 0x0f)
	gb.mainMemory.write(0xff85, 0xf0)
	gb.set16Reg(HL, 0xff85)
	gb.SUB_a_hl([1]uint8{0x96})
	assert.Equal(t, gb.get8Reg(A), uint8(0x1f))
	assert.Equal(t, gb.get8Reg(F), uint8(0x50)) // N_FLAG and C_FLAG is set
}

func TestSBC_a_r(t *testing.T) {
	gb := initGameboy()
	gb.SBC_a_r([1]uint8{0x98})
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0xc0)) // Z_FLAG and N_FLAG is set

	gb.set8Reg(A, 0x03)
	gb.set8Reg(B, 0x01)
	gb.set8Reg(F, 0x10)
	gb.SBC_a_r([1]uint8{0x98})
	assert.Equal(t, gb.get8Reg(A), uint8(0x01))
	assert.Equal(t, gb.get8Reg(F), uint8(0x40)) // N_FLAG is set

	gb.set8Reg(A, 0x87)
	gb.set8Reg(B, 0x0e)
	gb.set8Reg(F, 0x10)
	gb.SBC_a_r([1]uint8{0x98})
	assert.Equal(t, gb.get8Reg(A), uint8(0x78))
	assert.Equal(t, gb.get8Reg(F), uint8(0x60)) // N_FLAG and H_FLAG is set

	gb.set8Reg(A, 0x0f)
	gb.set8Reg(B, 0xef)
	gb.set8Reg(F, 0x10)
	gb.SBC_a_r([1]uint8{0x98})
	assert.Equal(t, gb.get8Reg(A), uint8(0x1f))
	assert.Equal(t, gb.get8Reg(F), uint8(0x50)) // N_FLAG and C_FLAG is set
}

func TestSBC_a_n(t *testing.T) {
	gb := initGameboy()
	gb.SBC_a_n([2]uint8{0xde, 0x00})
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0xc0)) // Z_FLAG and N_FLAG is set

	gb.set8Reg(A, 0x03)
	gb.set8Reg(F, 0x10)
	gb.SBC_a_n([2]uint8{0xde, 0x01})
	assert.Equal(t, gb.get8Reg(A), uint8(0x01))
	assert.Equal(t, gb.get8Reg(F), uint8(0x40)) // N_FLAG is set

	gb.set8Reg(A, 0x87)
	gb.set8Reg(F, 0x10)
	gb.SBC_a_n([2]uint8{0xde, 0x0e})
	assert.Equal(t, gb.get8Reg(A), uint8(0x78))
	assert.Equal(t, gb.get8Reg(F), uint8(0x60)) // N_FLAG and H_FLAG is set

	gb.set8Reg(A, 0x0f)
	gb.set8Reg(F, 0x10)
	gb.SBC_a_n([2]uint8{0xde, 0xef})
	assert.Equal(t, gb.get8Reg(A), uint8(0x1f))
	assert.Equal(t, gb.get8Reg(F), uint8(0x50)) // N_FLAG and C_FLAG is set
}

func TestSBC_a_hl(t *testing.T) {
	gb := initGameboy()
	gbROM := newGBROM()
	gb.mainMemory.cartridge = gbROM
	gb.SBC_a_hl([1]uint8{0x9e})
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0xc0)) // Z_FLAG and N_FLAG is set

	gb.set8Reg(A, 0x03)
	gb.set8Reg(F, 0x10)
	gb.mainMemory.write(0xff85, 0x01)
	gb.set16Reg(HL, 0xff85)
	gb.SBC_a_hl([1]uint8{0x9e})
	assert.Equal(t, gb.get8Reg(A), uint8(0x01))
	assert.Equal(t, gb.get8Reg(F), uint8(0x40)) // N_FLAG is set

	gb.set8Reg(A, 0x87)
	gb.set8Reg(F, 0x10)
	gb.mainMemory.write(0xff85, 0x0e)
	gb.set16Reg(HL, 0xff85)
	gb.SBC_a_hl([1]uint8{0x9e})
	assert.Equal(t, gb.get8Reg(A), uint8(0x78))
	assert.Equal(t, gb.get8Reg(F), uint8(0x60)) // N_FLAG and H_FLAG is set

	gb.set8Reg(A, 0x0f)
	gb.set8Reg(F, 0x10)
	gb.mainMemory.write(0xff85, 0xef)
	gb.set16Reg(HL, 0xff85)
	gb.SBC_a_hl([1]uint8{0x9e})
	assert.Equal(t, gb.get8Reg(A), uint8(0x1f))
	assert.Equal(t, gb.get8Reg(F), uint8(0x50)) // N_FLAG and C_FLAG is set
}

func TestAND_a_r(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0xcc)
	gb.set8Reg(B, 0xaf)
	gb.AND_a_r([1]uint8{0xa0}) // AND A B
	assert.Equal(t, gb.get8Reg(A), uint8(0x8c))
	assert.Equal(t, gb.get8Reg(F), uint8(0x20)) // H_FLAG is set

	gb.set8Reg(A, 0x12)
	gb.set8Reg(B, 0x00)
	gb.AND_a_r([1]uint8{0xa0}) // AND A B
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0xa0)) // H_FLAG and Z_FLAG are set
}

func TestAND_a_n(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0xcc)
	gb.AND_a_n([2]uint8{0xe6, 0xaf}) // AND A 0xaf
	assert.Equal(t, gb.get8Reg(A), uint8(0x8c))
	assert.Equal(t, gb.get8Reg(F), uint8(0x20)) // H_FLAG is set

	gb.set8Reg(A, 0x12)
	gb.set8Reg(B, 0x00)
	gb.AND_a_n([2]uint8{0xe6, 0x00}) // AND A 0x00
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0xa0)) // H_FLAG and Z_FLAG are set
}

func TestAND_a_hl(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0xcc)
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0xaf)
	gb.AND_a_hl([1]uint8{0xa6})
	assert.Equal(t, gb.get8Reg(A), uint8(0x8c))
	assert.Equal(t, gb.get8Reg(F), uint8(0x20)) // H_FLAG is set

	gb.set8Reg(A, 0x12)
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0x00)
	gb.AND_a_hl([1]uint8{0xa6})
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0xa0)) // H_FLAG and Z_FLAG are set
}

func TestOR_a_r(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0xaa)
	gb.set8Reg(B, 0x55)
	gb.OR_a_r([1]uint8{0xb0}) // OR b
	assert.Equal(t, gb.get8Reg(A), uint8(0xff))

	gb.set8Reg(A, 0x00)
	gb.set8Reg(B, 0x00)
	gb.OR_a_r([1]uint8{0xb0}) // OR b
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0x80)) // Z_FLAG is set
}

func TestOR_a_n(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0xaa)
	gb.OR_a_n([2]uint8{0xf6, 0x55}) // OR 0x55
	assert.Equal(t, gb.get8Reg(A), uint8(0xff))

	gb.set8Reg(A, 0x00)
	gb.OR_a_n([2]uint8{0xf6, 0x00}) // OR 0x00
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0x80)) // Z_FLAG is set
}

func TestOR_a_hl(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0xaa)
	gb.mainMemory.write(0xff85, 0x55)
	gb.set16Reg(HL, 0xff85)
	gb.OR_a_hl([1]uint8{0xb6}) // OR (HL)
	assert.Equal(t, gb.get8Reg(A), uint8(0xff))

	gb.set8Reg(A, 0x00)
	gb.mainMemory.write(0xff85, 0x00)
	gb.set16Reg(HL, 0xff85)
	gb.OR_a_hl([1]uint8{0xb6}) // OR (HL)
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0x80)) // Z_FLAG is set
}

func TestXOR_a_r(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0xaa)
	gb.set8Reg(B, 0xff)
	gb.XOR_a_r([1]uint8{0xa8}) // XOR B
	assert.Equal(t, gb.get8Reg(A), uint8(0x55))

	gb.set8Reg(A, 0xaa)
	gb.set8Reg(B, 0xaa)
	gb.XOR_a_r([1]uint8{0xa8}) // XOR B
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0x80)) // Z_FLAG is set
}

func TestXOR_a_n(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0xaa)
	gb.XOR_a_n([2]uint8{0xee, 0xff}) // XOR 0xff
	assert.Equal(t, gb.get8Reg(A), uint8(0x55))

	gb.set8Reg(A, 0xaa)
	gb.XOR_a_n([2]uint8{0xee, 0xaa}) // XOR 0xaa
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0x80)) // Z_FLAG is set
}

func TestXOR_a_hl(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0xaa)
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0xff)
	gb.XOR_a_hl([1]uint8{0xae}) // XOR (HL)
	assert.Equal(t, gb.get8Reg(A), uint8(0x55))

	gb.set8Reg(A, 0xaa)
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0xaa)
	gb.XOR_a_hl([1]uint8{0xae}) // XOR (HL)
	assert.Equal(t, gb.get8Reg(A), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0x80)) // Z_FLAG is set
}

func TestCP_a_r(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0x05)
	gb.set8Reg(B, 0x02)
	gb.CP_a_r([1]uint8{0xb8})                   // CP B
	assert.Equal(t, gb.get8Reg(F), uint8(0x60)) // H_FLAG and N_FLAG is set

	gb.set8Reg(A, 0x02)
	gb.set8Reg(B, 0x05)
	gb.CP_a_r([1]uint8{0xb8})                   // CP B
	assert.Equal(t, gb.get8Reg(F), uint8(0x50)) // C_FLAG and N_FLAG is set

	gb.set8Reg(A, 0x03)
	gb.set8Reg(B, 0x03)
	gb.CP_a_r([1]uint8{0xb8})                   // CP B
	assert.Equal(t, gb.get8Reg(F), uint8(0xc0)) // N_FLAG and Z_FLAG is set
}

func TestCP_a_n(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0x05)
	gb.set8Reg(B, 0x02)
	gb.CP_a_n([2]uint8{0xfe, 0x02})             // CP 0x02
	assert.Equal(t, gb.get8Reg(F), uint8(0x60)) // H_FLAG and N_FLAG is set

	gb.set8Reg(A, 0x02)
	gb.set8Reg(B, 0x05)
	gb.CP_a_n([2]uint8{0xfe, 0x05})             // CP 0x05
	assert.Equal(t, gb.get8Reg(F), uint8(0x50)) // C_FLAG and N_FLAG is set

	gb.set8Reg(A, 0x03)
	gb.set8Reg(B, 0x03)
	gb.CP_a_n([2]uint8{0xfe, 0x03})             // CP 0x03
	assert.Equal(t, gb.get8Reg(F), uint8(0xc0)) // N_FLAG and Z_FLAG is set
}

func TestCP_a_hl(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0x05)
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0x02)
	gb.CP_a_hl([1]uint8{0xbe})                  // CP (HL)
	assert.Equal(t, gb.get8Reg(F), uint8(0x60)) // H_FLAG and N_FLAG is set

	gb.set8Reg(A, 0x02)
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0x05)
	gb.CP_a_hl([1]uint8{0xbe})                  // CP (HL)
	assert.Equal(t, gb.get8Reg(F), uint8(0x50)) // C_FLAG and N_FLAG is set

	gb.set8Reg(A, 0x03)
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0x03)
	gb.CP_a_hl([1]uint8{0xbe})                  // CP (HL)
	assert.Equal(t, gb.get8Reg(F), uint8(0xc0)) // Z_FLAG and N_FLAG is set
}

func TestINC_r(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(B, 0x41)
	gb.INC_r([1]uint8{0x04})
	assert.Equal(t, gb.get8Reg(B), uint8(0x42)) // INC B
	assert.Equal(t, gb.get8Reg(F), uint8(0x00))

	gb.set8Reg(B, 0xff)
	gb.INC_r([1]uint8{0x04})
	assert.Equal(t, gb.get8Reg(B), uint8(0x00)) // INC B
	assert.Equal(t, gb.get8Reg(F), uint8(0xa0)) // H_FLAG and Z_FLAG is set

	gb.set8Reg(B, 0x0f)
	gb.INC_r([1]uint8{0x04})
	assert.Equal(t, gb.get8Reg(B), uint8(0x10)) // INC B
	assert.Equal(t, gb.get8Reg(F), uint8(0x20)) // H_FLAG is set
}

func TestINC_hl(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0x41)
	gb.INC_hl([1]uint8{0xbe}) // INC (HL)
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0x42))
	assert.Equal(t, gb.get8Reg(F), uint8(0x00))

	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0xff)
	gb.INC_hl([1]uint8{0xbe}) // INC (HL)
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0xa0)) // H_FLAG and Z_FLAG is set

	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0x0f)
	gb.INC_hl([1]uint8{0xbe}) // INC (HL)
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0x10))
	assert.Equal(t, gb.get8Reg(F), uint8(0x20)) // H_FLAG is set
}

func TestDEC_r(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(B, 0x41)
	gb.DEC_r([1]uint8{0x05}) // DEC B
	assert.Equal(t, gb.get8Reg(B), uint8(0x40))
	assert.Equal(t, gb.get8Reg(F), uint8(0x40)) // N_FLAG is set

	gb.set8Reg(B, 0x01)
	gb.DEC_r([1]uint8{0x05}) // DEC B
	assert.Equal(t, gb.get8Reg(B), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0xc0)) // N_FLAG and Z_FLAG is set

	gb.set8Reg(B, 0x00)
	gb.DEC_r([1]uint8{0x05}) // DEC B
	assert.Equal(t, gb.get8Reg(B), uint8(0xff))
	assert.Equal(t, gb.get8Reg(F), uint8(0x60)) // N_FLAG and H_FLAG is set
}

func TestDEC_hl(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0x41)
	gb.DEC_hl([1]uint8{0xbe}) // DEC (HL)
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0x40))
	assert.Equal(t, gb.get8Reg(F), uint8(0x40)) // N_FLAG is set

	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0x01)
	gb.DEC_hl([1]uint8{0xbe}) // DEC (HL)
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0x00))
	assert.Equal(t, gb.get8Reg(F), uint8(0xc0)) // N_FLAG and Z_FLAG is set

	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0x00)
	gb.DEC_hl([1]uint8{0xbe}) // DEC (HL)
	assert.Equal(t, gb.mainMemory.read(0xff85), uint8(0xff))
	assert.Equal(t, gb.get8Reg(F), uint8(0x60)) // M_FLAG and H_FLAG is set
}

func TestADD_hl_ss(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(HL, 0x0fff)
	gb.set16Reg(BC, 0x0001)
	gb.ADD_hl_ss([1]uint8{0x09})
	assert.Equal(t, gb.get16Reg(HL), uint16(0x1000))
	assert.Equal(t, gb.get8Reg(F), uint8(0x20))

	gb.set16Reg(HL, 0xf000)
	gb.set16Reg(BC, 0x1001)
	gb.ADD_hl_ss([1]uint8{0x09})
	assert.Equal(t, gb.get16Reg(HL), uint16(0x0001))
	assert.Equal(t, gb.get8Reg(F), uint8(0x10))
}

func TestADD_sp_e(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(SP, 0x0fff)
	gb.ADD_sp_e([2]uint8{0xe8, 0x01})
	assert.Equal(t, gb.get16Reg(SP), uint16(0x1000))
	assert.Equal(t, gb.get8Reg(F), uint8(0x20))

	gb.set16Reg(SP, 0xffff)
	gb.ADD_sp_e([2]uint8{0xe8, 0x01})
	assert.Equal(t, gb.get16Reg(SP), uint16(0x0000))
	assert.Equal(t, gb.get8Reg(F), uint8(0x30))
}

// TODO: Implement 16 bit INC test
// TODO: Implement 16 bit DEC test

func TestJP_nn(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(PC, 0x1234)
	gb.JP_nn([3]uint8{0xc3, 0x78, 0x56})
	assert.Equal(t, gb.get16Reg(PC), uint16(0x5678))
}

func TestJP_cc_nn(t *testing.T) {
	gb := initGameboy()
	gb.JP_cc_nn([3]uint8{0xc2, 0x34, 0x12})
	assert.Equal(t, gb.get16Reg(PC), uint16(0x1234))

	gb.set16Reg(PC, 0x0000)
	gb.modifyFlag(Z_FLAG, SET)
	gb.JP_cc_nn([3]uint8{0xc2, 0x34, 0x12})
	assert.Equal(t, gb.get16Reg(PC), uint16(0x0003))

	gb.set16Reg(PC, 0x0000)
	gb.modifyFlag(C_FLAG, SET)
	gb.JP_cc_nn([3]uint8{0xda, 0x34, 0x12})
	assert.Equal(t, gb.get16Reg(PC), uint16(0x1234))

	gb.set16Reg(PC, 0x0000)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.JP_cc_nn([3]uint8{0xda, 0x34, 0x12})
	assert.Equal(t, gb.get16Reg(PC), uint16(0x0003))
}

func TestJR_e(t *testing.T) {
	gb := initGameboy()
	gb.JR_e([2]uint8{0x18, 0x0a})
	assert.Equal(t, gb.get16Reg(PC), uint16(0x000c))

	gb.JR_e([2]uint8{0x18, 0xf9})
	assert.Equal(t, gb.get16Reg(PC), uint16(0x0007))
}

func TestJR_cc_e(t *testing.T) {
	gb := initGameboy()
	gb.JR_cc_e([2]uint8{0x20, 0x0a})
	assert.Equal(t, gb.get16Reg(PC), uint16(0x000c))

	gb.set16Reg(PC, 0x0000)
	gb.modifyFlag(Z_FLAG, SET)
	gb.JR_cc_e([2]uint8{0x20, 0xf9})
	assert.Equal(t, gb.get16Reg(PC), uint16(0x0002))

	gb.set16Reg(PC, 0x000c)
	gb.modifyFlag(C_FLAG, SET)
	gb.JR_cc_e([2]uint8{0x38, 0xf9})
	assert.Equal(t, gb.get16Reg(PC), uint16(0x0007))

	gb.set16Reg(PC, 0x000c)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.JR_cc_e([2]uint8{0x38, 0x34})
	assert.Equal(t, gb.get16Reg(PC), uint16(0x000e))
}

func TestJP_hl(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(HL, 0x1234)
	gb.JP_hl([1]uint8{0xe9})
	assert.Equal(t, gb.get16Reg(PC), uint16(0x1234))
}

func TestRLCA(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0x85)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.RLCA([1]uint8{0x00}) // RLCA
	// result should be 0x0b (which contradicts the GB programming manual)
	// discussion: https://hax.iimarckus.org/topic/1617/
	assert.Equal(t, uint8(0x0b), gb.get8Reg(A))
	assert.Equal(t, uint8(0x10), gb.get8Reg(F)) // C FLAG set, Z,H,N FLAG cleared
}

func TestRLA(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0x95)
	gb.modifyFlag(C_FLAG, SET)
	gb.RLA([1]uint8{0x00}) // RLA
	assert.Equal(t, uint8(0x2b), gb.get8Reg(A))
	assert.Equal(t, uint8(0x10), gb.get8Reg(F)) // C FLAG set, Z,H,N FLAG cleared
}

func TestRRCA(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0x3b)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.RRCA([1]uint8{0x00}) // RRCA
	assert.Equal(t, uint8(0x9d), gb.get8Reg(A))
	assert.Equal(t, uint8(0x10), gb.get8Reg(F)) // C FLAG set, Z,H,N FLAG cleared
}

func TestRRA(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(A, 0x81)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.RRA([1]uint8{0x00}) // RRA
	assert.Equal(t, uint8(0x40), gb.get8Reg(A))
	assert.Equal(t, uint8(0x10), gb.get8Reg(F)) // C FLAG set, Z,H,N FLAG cleared
}

func TestRLC_r(t *testing.T) {
	gb := initGameboy()
	gb.set8Reg(B, 0x85)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.RLC_r([1]uint8{0x05}) // RLC_r
	assert.Equal(t, uint8(0x0b), gb.get8Reg(B))
	assert.Equal(t, uint8(0x10), gb.get8Reg(F)) // C FLAG set, Z,H,N FLAG cleared
}

func TestRLC_hl(t *testing.T) {
	gb := initGameboy()
	gb.set16Reg(HL, 0xff85)
	gb.mainMemory.write(0xff85, 0x00)
	gb.modifyFlag(C_FLAG, CLEAR)
	gb.RLC_hl([1]uint8{0xbe}) // RLC_hl
	assert.Equal(t, uint8(0x00), gb.mainMemory.read(0xff85))
	assert.Equal(t, uint8(0x80), gb.get8Reg(F)) // Z FLAG set, C,H,N FLAG cleared
}
