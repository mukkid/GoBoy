package GoBoy

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestGet16Reg(t *testing.T) {
    regs := &Register{}
    assert.Equal(t, regs.get16Reg(BC), uint16(0x00))
}

func TestSet16Reg(t *testing.T) {
    regs := &Register{}
    regs.set16Reg(BC, 0xbeef)
    assert.Equal(t, regs.get16Reg(BC), uint16(0xbeef))
}

func TestGet8Reg(t *testing.T) {
    regs := &Register{}
    regs.set16Reg(BC, 0xfeed)
    assert.Equal(t, regs.get8Reg(C), uint8(0xed))
}

func TestSet8Reg(t *testing.T) {
    regs := &Register{}
    regs.set8Reg(A, 0xde)
    assert.Equal(t, regs.get8Reg(A), uint8(0xde))
}

func TestModifyFlag(t *testing.T) {
    regs := &Register{}
    regs.modifyFlag(Z_FLAG, 1)
    assert.Equal(t, regs.get8Reg(F), uint8(0x80))
    regs.modifyFlag(Z_FLAG, 0)
    assert.Equal(t, regs.get8Reg(F), uint8(0x00))
    regs.modifyFlag(N_FLAG, 1)
    regs.modifyFlag(H_FLAG, 1)
    regs.modifyFlag(C_FLAG, 1)
    assert.Equal(t, regs.get8Reg(F), uint8(0x70))
}
