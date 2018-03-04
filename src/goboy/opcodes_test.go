package GoBoy

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func initGameboy() *GameBoy {
    return &GameBoy{
        Register: &Register{},
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
    gb.mainMemory.rom0[0x1234] = 0x42
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
    gb.mainMemory.rom0[0x1234] = 0x42
    memValue := gb.mainMemory.read(0x1234)
    assert.Equal(t, memValue, uint8(0x42))
    gb.set16Reg(BC, 0x1234)
    gb.LD_a_bc([1]uint8{0x0a}) // LD A (BC)
    assert.Equal(t, gb.get8Reg(A), uint8(0x42))
}

func TestLD_a_de(t *testing.T) {
    gb := initGameboy()
    gb.mainMemory.rom0[0x1234] = 0x42
    memValue := gb.mainMemory.read(0x1234)
    assert.Equal(t, memValue, uint8(0x42))
    gb.set16Reg(DE, 0x1234)
    gb.LD_a_de([1]uint8{0x1a}) // LD A (DE)
    assert.Equal(t, gb.get8Reg(A), uint8(0x42))
}

func TestLD_a_nn(t *testing.T) {
    gb := initGameboy()
    gb.mainMemory.rom0[0x1234] = 0x42
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
    gb.set8Reg(A ,0x42)
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
    gb.mainMemory.rom0[0x0123] = 0x37
    gb.mainMemory.rom0[0x0124] = 0xa1
    gb.LD_hl_nn([3]uint8{0x2a, 0x23, 0x01}) // LD HL <- (0x0123)
    assert.Equal(t, gb.get16Reg(HL), uint16(0xa137))
}

func TestLD_dd_NN(t *testing.T) {
    gb := initGameboy()
    gb.mainMemory.rom0[0x0123] = 0xcd
    gb.mainMemory.rom0[0x0124] = 0xab
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
