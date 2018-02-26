package GoBoy

import "testing"

func TestLD_r_r(t *testing.T) {
    gb := &GameBoy{
        Register: &Register{},
    }
    gb.set8Reg(B, 0xfe)
    gb.LD_r_r(0x50) // LD D, B
    if gb.get8Reg(B) != 0xfe {
        t.Fatal()
    }
    if gb.get8Reg(D) != 0xfe {
        t.Fatal()
    }
    if gb.get16Reg(PC) != 0x0001 {
        t.Fatal()
    }
}

func TestLD_r_n(t *testing.T) {
    gb := &GameBoy{
        Register: &Register{},
    }
    gb.LD_r_n(0xead) // LD C, 0xad
    if gb.get8Reg(C) != 0xad {
        t.Fatal()
    }
    if gb.get16Reg(PC) != 0x0002 {
        t.Fatal()
    }
}
