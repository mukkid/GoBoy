package GoBoy

import "testing"

func TestGet16Reg(t *testing.T) {
    regs := &Register{}
    if regs.get16Reg(BC) != 0x00 {
        t.Fatal()
    }
}

func TestSet16Reg(t *testing.T) {
    regs := &Register{}
    regs.set16Reg(BC, 0xbeef)
    if regs.get16Reg(BC) != 0xbeef {
        t.Fatal()
    }
}

func TestGet8Reg(t *testing.T) {
    regs := &Register{}
    regs.set16Reg(BC, 0xfeed)
    if regs.get8Reg(C) != 0xed {
        t.Fatal()
    }
}

func TestSet8Reg(t *testing.T) {
    regs := &Register{}
    regs.set8Reg(A, 0xde)
    if regs.get8Reg(A) != 0xde {
        t.Fatal()
    }
}

func TestIncrementPc(t *testing.T) {
    regs := &Register{}
    regs.incrementPc(1)
    if regs.get16Reg(PC) != 0x0001 {
        t.Fatal()
    }
    regs.incrementPc(3)
    if regs.get16Reg(PC) != 0x0004 {
        t.Fatal()
    }
    regs.incrementPc(-2)
    if regs.get16Reg(PC) != 0x0002 {
        t.Fatal()
    }
}
