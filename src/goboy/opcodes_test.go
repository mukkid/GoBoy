package GoBoy

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func initGameboy() *GameBoy {
    return &GameBoy{Register: &Register{}}
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

func TestLD_sp_hl(t *testing.T) {
    gb := initGameboy()
    gb.set16Reg(HL, 0x1234)
    gb.LD_sp_hl([1]uint8{0xf9}) // LD SP HL (HL = 0x1234)
    assert.Equal(t, gb.get16Reg(SP), uint16(0x1234))
}
